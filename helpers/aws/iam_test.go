package aws_helpers

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
		project_name := "test-project"
		region := "us-east-1"
		accountID := "123456789012"

		iamUser, err := createIAMUserResource(ctx, project_name, region, accountID)
		assert.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)

		// Check if the IAM User has the correct tags
		pulumi.All(iamUser.URN(), iamUser.Tags).ApplyT(func(all []interface{}) error {
			defer wg.Done() // Ensure Done is called
			fmt.Printf("All: %v\n", all)

			urn := all[0].(pulumi.URN)
			tags := all[1].(map[string]string)
			nameTag, ok := tags["Name"]
			assert.True(t, ok, "Missing Name tag on User %v", urn)
			assert.Equal(t, project_name, nameTag, "Incorrect Name tag on User %v", urn)



			return nil
		})

		wg.Wait()
		return nil
	}, pulumi.WithMocks("project", "stack", mocks.Mocks(0)))
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

