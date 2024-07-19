package aws

import (
	"maos-cloud-project-api/mocks"
	"sync"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestCreateIAMUser(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		iamUser, err := createIAMUserResource(ctx, "project_name", "region", "account_id")
		assert.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(3)

		// Check if the IAM User has the correct tags
		pulumi.All(iamUser.arn, iamUser.user.Arn).ApplyT(func ( all []interface{}) error {
			urn := all[0].(pulumi.URN)
			tags := all[1].(map[string]string)
			assert.Containsf(t, tags, "Name", "Missing Name tag on User %v", urn)
			wg.Done()
			return nil
			
		})

		// Check if the IAM User has the correct s3 policy attached
		pulumi.All(iamUser.arn, iamUser.user.Arn).ApplyT(func ( all []interface{}) error {
			urn := all[0].(pulumi.URN)
			policy := all[1].(string)
			assert.Containsf(t, policy, "s3Policy", "Missing S3 policy on User %v", urn)
			wg.Done()
			return nil
		})

		// Check if the IAM User has the correct ECR policy attached
		pulumi.All(iamUser.arn, iamUser.user.Arn).ApplyT(func ( all []interface{}) error {
			urn := all[0].(pulumi.URN)	
			policy := all[1].(string)
			assert.Containsf(t, policy, "ecrPolicy", "Missing ECR policy on User %v", urn)
			wg.Done()
			return nil
		})

		wg.Wait()
		return nil
	}, pulumi.WithMocks("project_name", "account_id", mocks.Mocks(0)))
	assert.NoError(t, err)
}

