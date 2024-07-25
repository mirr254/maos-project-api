package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"maos-cloud-project-api/controllers/aws"
	"maos-cloud-project-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/sirupsen/logrus"
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

	projectConfig, err := utils.BuildProjectConfig(project_name, region)
	if err != nil {
		logrus.Error("Could not build project config", err)
		return auto.UpdateSummary{}, err
	}


    projectConfigMap := utils.ConvertProjectConfigToAutoConfigMap(projectConfig)
	stck.SetAllConfig(ctx, projectConfigMap)

	//deploy stack
	upRes, err := stck.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		logrus.Error("Could not deploy stack", err)
		return auto.UpdateSummary{}, err
	}

	// TODO: WHEN the stack is created successfully, redirect to dashboard and run
	// CreateIAMUser function in the background

		// Generate Pulumi.yaml file dynamically
	err = utils.GeneratePulumiYAML(projectConfig, fmt.Sprintf("Pulumi.%s.yaml", stack_name))
	if err != nil {
		fmt.Println("Could not generate Pulumi.yaml file: ", err)
		return auto.UpdateSummary{}, err
	}

	err = aws.CreateIAMUser(project_name, region, stack_name)
	if err != nil {
		logrus.Error("Could not create IAM user", err)
		return auto.UpdateSummary{}, err
	}

	return upRes.Summary, nil
	
}

//Deletes a particular stack and all the associated resources
func DeleteStack(c *gin.Context) {
	projectName := c.Param("project_name")
	stackName   := c.Param("stack_name")
    
	ctx := context.Background()

	s, err := auto.SelectStackInlineSource(ctx, stackName, projectName , PulumiProgram, auto.Program(PulumiProgram))
	if err != nil {
		// check if stack exists
		if auto.IsSelectStack404Error(err){
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack doesn't exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf(err.Error())})
		return
	}

	// destroy stack
	_, err = s.Destroy(ctx, optdestroy.ProgressStreams(os.Stdout))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not destroy stack"})
		return
	}

	//delete the stack and all associated history and config
	err = s.Workspace().RemoveStack(ctx, stackName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove stack"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stack deleted successfully"})
	
}

