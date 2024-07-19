package controllers

import (
	"fmt"
	"math/rand"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/sirupsen/logrus"
)

type Project struct {
	ProjectName   string `json:"project_name"`
	Region        string `json:"region"`
	StackName     string `json:"stack_name"`
	CloudProvider string `json:"cloud_provider"`
}

type ProjectResponse struct {
	StackName     string `json:"stack_name"`
	URL           string `json:"url"`
	CloudProvider string `json:"cloud_provider"`
	// Status      string `json:"status"`
}

/*
CreateProject: function is responsible for creating a new project on pulumi dashboard.
A Client can have multiple projects.

*/

func CreateProject(c *gin.Context) {

	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	suffixedProjectName := suffixProjectName(project.ProjectName)

	UpdateSummary, err := CreateStack(suffixedProjectName, project.Region, project.StackName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create stack"})
		return
	}

	updateSummary := UpdateSummary
	outputJson, err := json.MarshalIndent(updateSummary, "", "  ")
	if err != nil {
		logrus.Error("Failed to marshal update summary: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not marshal output"})
		return
	}

	logrus.Info("Update Summary: ", string(outputJson))


	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully", "project_name": suffixedProjectName, "stack_name": project.StackName})
	return

}

func PulumiProgram(ctx *pulumi.Context) error {
	// Implement your Pulumi program here
	ctx.Export("Success!", pulumi.Sprintf("success"))
	return nil
}

/*
This function adds a suffix to the provided project name to avoid duplicate names
Params: projectName
Returns: suffix-projectname
*/
func suffixProjectName(projectName string) string {

	rand.Seed(time.Now().UnixNano())
	min := 100
	max := 10000

	fmt.Sprintf("Project name is %s", projectName+"-"+strconv.Itoa(rand.Intn(max-min+1)))

	return projectName + "-" + strconv.Itoa(rand.Intn(max-min+1))
}
