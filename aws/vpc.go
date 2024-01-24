package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	// "io/ioutil"
	// "strings"

	"github.com/gorilla/mux"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//define request/response types various rest ops
type Vpc struct {
	VpcName      string `json:"vpc_name"`
	Tags         VpcTags `json:"vpc_tags"`
}

type VpcTags struct {
	VpcName      string `json:"vpc_name"`
	ProjectName  string `json:"project_name"`
	StackName    string `json:"environment"`
	CidrBlock    string `json:"cidr_block"`
	Region       string `json:"region"`
	
}

type VpcResponse struct {
	VpcID       string `json:"vpc_id"`
	Status      string `json:status`
}

func CreateAwsVPC(w http.ResponseWriter, req *http.Request) {
    
    w.Header().Set("Content-type", "application/json")

	params := mux.Vars(req)
	projectName = params["project_name"]
	stackName = params["stack_name"]

    var vpc Vpc
	err := json.NewDecoder(req.Body).Decode(&vpc)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Falied to parse the VPC args ")
		return 
	}
    // for debugging
	vpc_req_data, _ := json.MarshalIndent( vpc, "" ,"\t" )
	fmt.Println(string(vpc_req_data))

	ctx := context.Background()

	tags     := vpc.Tags
	program  := createVpcProgram( tags )

	//aresociate with a stack
	s, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, program)
	if err != nil {
		// if stack doesn't exists return 404 exit
		if auto.IsSelectStack404Error(err) {
			w.WriteHeader(404)
			fmt.Fprintf(w, fmt.Sprintf("Stack %v for the project %v doesn't exists", stackName, projectName))
			return
		}

		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return 

	}
	
	s.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: tags.Region })
	
	// deploy stack
	// write all update logs to stdout so we can see the updates
	upVpcRes, err := s.Up(ctx, optup.ProgressStreams(os.Stdout))

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return 
	}

	// Subnet setup
	subnetDetails := make(map[string]interface{})
	privateSubnet := make(map[string]interface{})
	publicSubnet  := make(map[string]interface{})

	privateSubnet["vpc_id"] = upVpcRes.Outputs["vpc_id"].Value.(string)
	privateSubnet["subnet_cidr_block"] = "10.0.1.0/16"
	// privateSubnet["availability_zone"] = 

	subnetDetails["private_subnet"] = privateSubnet
	subnetDetails["public_subnet"] = publicSubnet

	//Create the subents here TODO: FIX ME
 
	// vpcID, err := s.GetConfig(ctx, "vpc_id")
	response := &VpcResponse{
		VpcID: upVpcRes.Outputs["vpc_id"].Value.(string),
	}

	
	json.NewEncoder(w).Encode(&response)

}

//this defines the our pulumi VPC in terms of the contents that the caller parses in
// allows dynamically deploy VPC based on user defined values from the REST body

func createVpcProgram( tags_ VpcTags ) pulumi.RunFunc{
	return func(ctx *pulumi.Context) error {

        vpc_tags := make(map[string]string)

		vpc_tags["region"]     = tags_.Region
		vpc_tags["cidr_block"] = tags_.CidrBlock 
		vpc_tags["project_name"] = tags_.ProjectName
        vpc_tags["environment"]  = tags_.StackName 
		vpc_tags["vpc_name"]     = tags_.VpcName

		 // for debugging
		// var tagsvpc VpcTags
		data, _ := json.Marshal(tags_)
		_ = json.Unmarshal(data, &vpc_tags)
		vpc_req_data, _ := json.MarshalIndent( vpc_tags, "" ,"\t" )

		fmt.Println(string(vpc_req_data))

		

		vpcArgs := &ec2.VpcArgs{
			CidrBlock: pulumi.String(tags_.CidrBlock),
			Tags:     pulumi.ToStringMap(vpc_tags),
		}
		vpc, err := ec2.NewVpc(ctx, tags_.VpcName, vpcArgs) 
		if err != nil {
			fmt.Printf("An error occured while creating the vpc %v \n", err.Error())
			return err
		}

		
		
		ctx.Export("vpc_id", vpc.ID() )

		return nil
	}
}

// delete a vpc
func DeleteVPC( w http.ResponseWriter, req *http.Request ) {
	w.Header().Set("Content-Type", "application/json")

	ctx := context.Background()
	params := mux.Vars(req)
	vpcID := params["vpc_id"]

	stackName = params["stack_name"]
	projectName = params["project_name"]

	s, err := auto.SelectStackInlineSource(ctx, stackName, projectName, CreateEmptyProgram())
	if err != nil {
		//if stack doesn't exist, 404
		if auto.IsSelectStack404Error(err) {
			w.WriteHeader(404)
			fmt.Fprintf(w, fmt.Sprintf("stack %q not found", stackName))
			return 
		}

		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
		return 
	}
    
	dep, _ := s.Export(ctx)
    // import/export is backwards compatible, and we must write code specific to the verison we're dealing with.
	if dep.Version != 3 {
		panic("expected deployment version 3")
	}
	var depState apitype.DeploymentV3

    err = json.Unmarshal(dep.Deployment, &depState)
    if err != nil {
		print(err)
		return
	}
	//debug purposes
	// resources, _ := json.MarshalIndent(depState, "", "\t")
	// fmt.Println(string(resources))

	depState.Resources = RemoveResource(depState, depState.Resources, vpcID )
	bytes, err := json.Marshal(depState)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, fmt.Sprintf("Couldn't Marshal"))
	}

	dep.Deployment = bytes
	//import edited deployment status back to the stack
	err = s.Import(ctx, dep)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	fmt.Printf("Done\n")
	
	// newResources, _ := json.MarshalIndent(depState, "", "\t")
	// fmt.Println(string(newResources))
   
	delResponse := &VpcResponse{
		VpcID: vpcID,
		Status: "Deleted",
	} 

	json.NewEncoder(w).Encode(&delResponse)
}

func UpdateVpc( w http.ResponseWriter, req *http.Request ) {

    w.Header().Set("Content-Type",  "application/json")
	var vpcUpdateReq Vpc
	err := json.NewDecoder(req.Body).Decode(&vpcUpdateReq)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "failed to parse request")
		return
	}

	ctx := context.Background()

	params := mux.Vars(req)
	stackName := params["stack_name"]
	tagsToEdit := vpcUpdateReq.Tags

	s, err := auto.SelectStackInlineSource(ctx, stackName, projectName, CreateEmptyProgram())
	if err != nil {
		if auto.IsSelectStack404Error(err) {
			w.WriteHeader(404)
			fmt.Fprintf(w, fmt.Sprintf("Stack %q not found", stackName))
			return
		}

		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
		return
	}

	// dep, err := s.Export(ctx)
    // fmt.Printf("Deployment numerb: %v", dep.Version)

	s.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "us-east-1"})
	s.SetConfig(ctx, "aws:vpc_name", auto.ConfigValue{Value: tagsToEdit.VpcName })

	//deploy 
	//write all of the updates logs to stdout so we can watch progreres
	upRes, err := s.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		//check if there is another update in progreres, return 409
		if auto.IsConcurrentUpdateError(err){
			w.WriteHeader(409)
			fmt.Fprintf(w, "There is an update on stack %q in progreres", stackName)
			return
		}

		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
		return
	}
   
	response := &VpcResponse{
		VpcID: upRes.Outputs["vpc_id"].Value.(string),
	}
    
	json.NewEncoder(w).Encode(&response)

}
