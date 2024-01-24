package aws

import (
	"fmt"
	"github.com/pulumi/pulumi/sdk/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var projectName string
var stackName string

func CreateEmptyProgram() pulumi.RunFunc{
	return func(ctx *pulumi.Context) error {

		return nil 
	}
}

//resourceID is the provider-assigned resource, if any, for custom resources.
//depState is the current deployment state, includes secrets configs, cofigMaps, resources in a particular stack
func RemoveResource( depState apitype.DeploymentV3, resSlice []apitype.ResourceV3, resourceID string ) ( []apitype.ResourceV3) {
	
	for index, resource := range depState.Resources {
        if resource.ID.String() == resourceID  {
			fmt.Printf("Resource found. Deleting...\n")
            return append(depState.Resources[0:index], depState.Resources[index+1:]...)
		}
	}

	return resSlice

}
