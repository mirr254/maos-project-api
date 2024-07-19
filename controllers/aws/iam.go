package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/sirupsen/logrus"
)

type iamUser struct {
	arn   pulumi.StringOutput
	user  *iam.User
}

// CreateIAMUser: function is responsible for creating a new IAM user in the stack.

func CreateIAMUser(project_name, region, stack_name string) (auto.UpdateSummary, error) {

	logrus.Info("Creating IAM user: ", project_name)

	ctx := context.Background()

	// a closure that captures the params and calls the program function
	programWithParams := func( ctx *pulumi.Context) (error) {

		_, err := createIAMUserResource(ctx, project_name, region, stack_name)
		return err
	}
	// Create a new local Pulumi workspace with a specified project
	localWorkspace, err := auto.NewLocalWorkspace(ctx, auto.Program(programWithParams), auto.Project(workspace.Project{
		Name:    tokens.PackageName(project_name),
		Runtime: workspace.NewProjectRuntimeInfo("go", nil),
	}))
	if err != nil {
		logrus.Error("Could not create workspace: ", err)
		return auto.UpdateSummary{}, err
	}

	newStack, err := auto.SelectStack(ctx, stack_name, localWorkspace)
	if err != nil {
		logrus.Error("Could not select stack: ", err)
		return auto.UpdateSummary{}, err
	}

	//deploy
	upRes, err := newStack.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		logrus.Error("Could not deploy stack: ", err)
		return auto.UpdateSummary{}, err
	}

	logrus.Info("Stack deployed successfully: ", upRes.Summary)
	return upRes.Summary, nil

}

// func to to create a new IAM user

func createIAMUserResource(ctx *pulumi.Context, project_name, region, account_id string) (*iamUser, error) {
	// Create an IAM user
	user, err := iam.NewUser(ctx, project_name, &iam.UserArgs{
		Name: pulumi.String( project_name ),
		Tags: pulumi.StringMap{
			"Name": pulumi.String( project_name ),
		},
	})
	if err != nil {
		return nil, err
	}

	// Define the S3 policy
	s3Policy := s3Policy(project_name)

	// TODO: Remove Hardcoded account_id
	ecrPolicy := ecrPolicy(region, account_id, project_name) 
	_, err = iam.NewUserPolicy(ctx, "s3Policy", &iam.UserPolicyArgs{
		User:   user.Name,
		Policy: pulumi.String(s3Policy),
	})
	if err != nil {
		return nil, err
	}

	// Attach the ECR policy to the user
	_, err = iam.NewUserPolicy(ctx, "ecrPolicy", &iam.UserPolicyArgs{
		User:   user.Name,
		Policy: pulumi.String(ecrPolicy),
	})
	if err != nil {
		return nil, err
	}

	return &iamUser{
		user: user,
		arn:  user.Arn,
	} ,nil
}

// function that return ecr policy
func ecrPolicy(region, account_id, project_name string) string {
	return fmt.Sprintf( `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"ecr:GetDownloadUrlForLayer",
					"ecr:BatchGetImage",
					"ecr:BatchCheckLayerAvailability",
					"ecr:PutImage",
					"ecr:InitiateLayerUpload",
					"ecr:UploadLayerPart",
					"ecr:CompleteLayerUpload"
				],
				"Resource": [
					"arn:aws:ecr:%s:%s:repository/%s"
				]
			}
		]
	}`, region, account_id, project_name)
}

// S3 policy
func s3Policy(project_name string) string {
	return fmt.Sprintf( `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"s3:ListBucket",
					"s3:GetObject",
					"s3:PutObject"
				],
				"Resource": [
					"arn:aws:s3:::%s",
					"arn:aws:s3:::%s/*"
				]
			}
		]
	}`, project_name, project_name)
}