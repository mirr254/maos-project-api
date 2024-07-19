package aws

import (
	"fmt"
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
		wg.Add(1)

		// Check if the IAM User has the correct tags
		pulumi.All(iamUser.arn, iamUser.user.Arn).ApplyT(func ( all []interface{}) error {

			fmt.Printf("All: %v\n", all)

			// arn := all[0].(string)
			// tags := all[1].(map[string]interface{})
			// assert.Containsf(t, tags, "Name", "Missing Name tag on User %v", arn)
			wg.Done()
			return nil
			
		})

		// Check if the IAM User has the correct s3 policy attached
		// pulumi.All(iamUser.arn, iamUser.user.Arn).ApplyT(func ( all []interface{}) error {
		// 	urn := all[0].(pulumi.URN)
		// 	policy := all[1].(string)
		// 	assert.Containsf(t, policy, "s3Policy", "Missing S3 policy on User %v", urn)
		// 	wg.Done()
		// 	return nil
		// })

		// // Check if the IAM User has the correct ECR policy attached
		// pulumi.All(iamUser.arn, iamUser.user.Arn).ApplyT(func ( all []interface{}) error {
		// 	urn := all[0].(string)	
		// 	policy := all[1].(string)
		// 	assert.Containsf(t, policy, "ecrPolicy", "Missing ECR policy on User %v", urn)
		// 	wg.Done()
		// 	return nil
		// })

		wg.Wait()
		return nil
	}, pulumi.WithMocks("project_name", "account_id", mocks.Mocks(0)))
	assert.NoError(t, err)
}

func TestECRPolicy(t *testing.T) {
	region := "us-east-1"
	account_id := "123456789012"
	project_name := "project_name"
	policy := ecrPolicy(region, account_id, project_name)
	assert.Contains(t, policy, "ecr:GetDownloadUrlForLayer")
	assert.Contains(t, policy, "ecr:BatchGetImage")
	assert.Contains(t, policy, "ecr:BatchCheckLayerAvailability")
	assert.Contains(t, policy, "ecr:PutImage")
	assert.Contains(t, policy, "ecr:InitiateLayerUpload")
	assert.Contains(t, policy, "ecr:UploadLayerPart")
	assert.Contains(t, policy, "ecr:CompleteLayerUpload")
}

func TestS3Policy(t *testing.T) {
	project_name := "project_name"
	policy := s3Policy(project_name)
	assert.Contains(t, policy, "s3:GetObject")
	assert.Contains(t, policy, "s3:ListBucket")
	assert.Contains(t, policy, "s3:PutObject")
}

