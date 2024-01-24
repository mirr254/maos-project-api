package aws

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/gorilla/mux"
// 	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
// 	"github.com/pulumi/pulumi/sdk/v3/go/auto"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )

// type Subnet struct {
// 	Name      string `json:"subnet_name"`
// 	VpcId     string `json:"vpc_id"`
// 	CidrBlock string `json:"subnet_cidr_block"`
// 	AvailabilityZone string `json:"availability_zone"`
// 	Tags      SubnetTags `json:"subnet_tags`
// }

// type SubnetTags struct {
// 	Name             string `json:"subnet_name"`
// 	ProjectName      string `json:"project_name"`
// 	VpcId            string `json:"vpc_id"`
// 	StackName        string `json:"environment"`
// 	AvailabilityZone string `json:"availability_zone"`
// 	CidrBlock        string `json:"cidr_block"`
// }

// type SubnetReponse struct {
// 	SubnetID      string `json:"subnet_id"`
// }

// func CreateSubnetProgram( tags_ SubnetTags ) pulumi.RunFunc{
// 	return func(ctx *pulumi.Context) error {

//         subnet_tags := make(map[string]string)

// 		subnet_tags["cidr_block"]   = tags_.CidrBlock 
// 		subnet_tags["project_name"] = tags_.ProjectName
//         subnet_tags["environment"]  = tags_.StackName 
// 		subnet_tags["subnet_name"]  = tags_.Name

// 		 // for debugging
// 		// var tagsvpc VpcTags
// 		data, _ := json.Marshal(tags_)
// 		_ = json.Unmarshal(data, &subnet_tags)
// 		subnet_req_data, _ := json.MarshalIndent( subnet_tags, "" ,"\t" )

// 		fmt.Println(string(subnet_req_data))

		
        
// 		subnetArgs := &ec2.SubnetArgs{
// 			CidrBlock: pulumi.String(tags_.CidrBlock),
// 			VpcId: pulumi.String(tags_.VpcId),
// 			Tags:     pulumi.ToStringMap(subnet_tags),
// 		}
// 		subnet, err := ec2.NewSubnet(ctx, tags_.Name, subnetArgs) 
// 		if err != nil {
// 			fmt.Printf("An error occured while creating the subnet %v \n", err.Error())
// 			return err
// 		}

				
// 		ctx.Export("subnet_id", subnet.ID() )

// 		return nil
// 	}
// }


// func CreateSubnets ( subnetDetails map[string]interface{} ) {
	
// 	projectName := subnetDetails[ "projectName" ]
// 	stackName   := subnetDetails[ "stack_name" ]

// 	//subent details map to string
// 	subnetDetails, _ = json.Marshal(subnetDetails)

// 	var subnet Subnet
// 	err := json.NewDecoder(subnet_details).Decode(&subnet)
// 	if err != nil {
// 		w.WriteHeader(400)
// 		fmt.Fprint(w, "Failed to parse the subnet args")
// 		return
// 	}

// 	ctx := context.Background()

// 	tags  := subnet.Tags
// 	program := createSubnetProgram( tags )

// 	s, err := auto.UpsertStackInlineSource(ctx, string(stackName), projectName, program)
// 	if err != nil {
// 		//if stack doesn't exist 404
// 		if auto.IsSelectStack404Error(err) {
// 			w.WriteHeader(404)
// 			fmt.Fprintf(w, fmt.Sprintf("Stack %v for the project %v doesn't exists", stackName, projectName))
// 			return 
// 		}

// 		w.WriteHeader(500)
// 		fmt.Fprint(w, err.Error())
// 		return 
// 	}

// 	cfg := auto.ConfigMap{
// 		"aws:subnet:name": {Value: tags.Name},
// 		"aws:subnet:vpc_id": { Value: tags.VpcId },
// 		"aws:subnet:az": { Value: tags.AvailabilityZone },
// 		"aws:subent:cidrblock": { Value: tags.CidrBlock },
// 		// "vpc_id":    {Value: "encrypted", Secret: true},
// 	}

// 	s.SetAllConfig(ctx, cfg)


// }
