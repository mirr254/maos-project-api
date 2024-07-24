package mocks

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Mocks int

func (Mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	pulumi.Printf("TOKEN TYPE: ", args.TypeToken)
	outputs := args.Inputs.Mappable()
	if args.TypeToken == "aws:iam/user:User" {
		outputs["arn"] = resource.NewStringProperty("arn:aws:iam::123456789012:user/" + args.Name)
		outputs["id"] = resource.NewStringProperty("123456789012")
	}
	//check vpc 
	// if args.TypeToken == "aws:ec2/vpc:Vpc" {
	// 	outputs["vpc_id"] = resource.NewStringProperty("vpc-1234567890")
	// 	outputs["cidr_block"] = resource.NewStringProperty("10.10.10.10/16")
	// 	outputs["tags"] = resource.NewPropertyMapFromMap(map[string]interface{}{
	// 		"Name": args.Name,
	// 	})
	// }

	return args.Name + "_id", resource.NewPropertyMapFromMap(outputs), nil
}

func (Mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	outputs := map[string]interface{}{}
	return resource.NewPropertyMapFromMap(outputs), nil
}
