package aws

import (
	"fmt"
	"maos-cloud-project-api/mocks"
	"sync"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	awsec2 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/stretchr/testify/assert"
)

// func TestCreateVPC(t *testing.T) {
// 	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
// 		project_name := "test-project"
// 		vpc, err := CreateVPCResource(ctx, project_name)
// 		assert.NoError(t, err)

// 		var wg sync.WaitGroup
// 		wg.Add(3)

// 		// Check if subnets greater than 3
// 		pulumi.All(vpc.vpc.URN(), vpc.vpc.Subnets).ApplyT(func(all []interface{}) error {
// 			defer wg.Done()
// 			fmt.Printf("VPC All: %v\n", all)

// 			urn := all[0].(pulumi.URN)
// 			subnets := all[1].([]awsec2.Subnet)
// 			assert.GreaterOrEqual(t, len(subnets), 3, "VPC %v should have at least 3 subnets", urn)

// 			wg.Done()
// 			return nil
// 		})
// 		wg.Wait()
// 		return nil
// 	}, pulumi.WithMocks("project", "stack", mocks.Mocks(0)))
// 	assert.NoError(t, err)
// }

func TestSecurityGroup(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		project_name := "test-project"
		vpc, err := CreateVPCResource(ctx, project_name)
		assert.NoError(t, err)

		// Create a new Security Group
		securityGroup, err := createSecurityGroup(ctx, project_name, vpc)
		assert.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)

		// Check if port 22 is open
		pulumi.All(securityGroup.group.URN(), securityGroup.group.Ingress).ApplyT(func(all []interface{}) error {
			
			fmt.Printf("Security Group All: %v\n", all)

			urn := all[0].(pulumi.URN)
			ingress := all[1].([]awsec2.SecurityGroupIngress)

			for _, i := range ingress {
				openToInternet := false 
				for _, cidr := range i.CidrBlocks {
					if cidr == "0.0.0.0/0" {
						openToInternet = true
						break
					}
			}
			assert.Falsef(t, i.FromPort == 22 && openToInternet, "Port 22 is not open to the internet (CIDR 0.0.0.0/0) on Security Group %v", urn)

		   }
			wg.Done()
			return nil
		})

		wg.Wait()
		return nil
	}, pulumi.WithMocks("project", "stack", mocks.Mocks(0)))
	assert.NoError(t, err)
}
