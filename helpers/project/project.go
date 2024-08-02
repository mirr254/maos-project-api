package project_helpers

import (
	"context"
	"fmt"
	"os"

	aws_helpers "maos-cloud-project-api/helpers/aws"
	"maos-cloud-project-api/utils"

	"github.com/sirupsen/logrus"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	
)

/*
Creates a new stack for a particular project. Here we understand stack at environment.
Prod, Staging, Dev
Return the stack name if created successfully

*/
func CreateStack( project_name, region, stack_name string) (auto.UpdateSummary, error) {

	ctx := context.Background()

	stck, err := auto.NewStackInlineSource(ctx, stack_name, project_name, PulumiProgram, auto.Program(PulumiProgram) )
	if err != nil {
		if auto.IsCreateStack409Error(err) {
			logrus.Error("Stack Exists error: ", err)
			return auto.UpdateSummary{}, err
		} else {
			logrus.Error("Stack Error ", err)
			return auto.UpdateSummary{}, err
		}

	}

	projectConfig, err := utils.BuildProjectConfig(project_name, region, "aws")
	if err != nil {
		logrus.Error("Could not build project config", err)
		return auto.UpdateSummary{}, err
	}

	// Generate Pulumi.yaml file dynamically
	err = utils.GeneratePulumiYAML(projectConfig, fmt.Sprintf("Pulumi.%s.yaml", stack_name))
	if err != nil {
		fmt.Println("Could not generate Pulumi.yaml file: ", err)
		return auto.UpdateSummary{}, err
	}


    projectConfigMap := utils.ConvertProjectConfigToAutoConfigMap(projectConfig)
	for key, value := range projectConfigMap {
		err := stck.SetConfig(ctx, key, value)
		if err != nil {
			logrus.Error("Could not set config ", err)
			return auto.UpdateSummary{}, err
		}
	}

	logrus.Info("Config set successfully")

	//deploy stack
	upRes, err := stck.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		logrus.Error("Could not deploy stack", err)
		return auto.UpdateSummary{}, err
	}

	// TODO: WHEN the stack is created successfully, redirect to dashboard and run
	// CreateIAMUser function in the background

	

	err = aws_helpers.CreateIAMUser(project_name, region, stack_name)
	if err != nil {
		logrus.Error("Could not create IAM user", err)
		return auto.UpdateSummary{}, err
	}

	return upRes.Summary, nil
	
}

func PulumiProgram(ctx *pulumi.Context) error {
	// Implement your Pulumi program here
	ctx.Export("Success!", pulumi.Sprintf("success"))
	return nil
}