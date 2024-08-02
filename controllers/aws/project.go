package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	project_helpers "maos-cloud-project-api/helpers/project"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	
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

	UpdateSummary, err := project_helpers.CreateStack(suffixedProjectName, project.Region, project.StackName)
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


//Deletes a particular stack and all the associated resources
func DeleteStack(c *gin.Context) {
	projectName := c.Param("project_name")
	stackName   := c.Param("stack_name")
    
	ctx := context.Background()

	s, err := auto.SelectStackInlineSource(ctx, stackName, projectName , project_helpers.PulumiProgram, auto.Program(project_helpers.PulumiProgram))
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
