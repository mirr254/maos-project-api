package mocks

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Mocks int

func (Mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	pulumi.Printf("TOKEN TYPE: %s\n\n", args.TypeToken)
	outputs := args.Inputs.Mappable()
	switch args.TypeToken {
	case "aws:iam/user:User":
		outputs["arn"] = pulumi.String("arn:aws:iam::123456789012:user/" + args.Name)
		outputs["id"] = pulumi.String("123456789012")
		outputs["tags"] = resource.NewObjectProperty(resource.PropertyMap{
			"Name": resource.NewStringProperty(args.Name),
		})

	case "aws:iam/userPolicy:UserPolicy":
		outputs["id"] = resource.NewStringProperty(args.Name + "_policy_id")
		outputs["name"] = resource.NewStringProperty(args.Name)
		outputs["policy"] = args.Inputs["policy"]
		outputs["user"] = args.Inputs["user"]
	}
	return args.Name + "_id", resource.NewPropertyMapFromMap(outputs), nil
}

func (Mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	outputs := map[string]interface{}{}
	return resource.NewPropertyMapFromMap(outputs), nil
}
