package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"maos-cloud-project-api/utils"
)

type Stack struct {
	StackName   string `json:"stack_name"`
	ProjectName string `json:"project_name"`
	Region      string `json:"region"`
}

type StackResponse struct {
	StackName  string `json:"stack_name"`
}

/*
Creates a new stack for a particular project. Here we understand stack at environment.
Prod, Staging, Dev
Return the stack name if created successfully

*/
func CreateStack( c *gin.Context) {
    
	var stack Stack
    if err := c.ShouldBindJSON(&stack); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectConfig := utils.BuildProjectConfig(stack.ProjectName, stack.Region)

	// for debugging purposes
	stackData, _ := json.MarshalIndent(stack, "", "\t")
	fmt.Println(string(stackData))
    
	ctx := context.Background()

	stackName := stack.StackName
	ProjectName := stack.ProjectName
	s, err := auto.NewStackInlineSource(ctx, stackName, ProjectName, PulumiProgram, auto.Program(PulumiProgram) )
	if err != nil {
		if auto.IsCreateStack409Error(err) {
			logrus.Error("Stack Exists error: ", err)
			c.JSON(http.StatusConflict, gin.H{"error": "Stack already exists"})
			return
		} else {
			logrus.Error("Stack Error ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create stack"})
			return
		}

	}


    projectConfigMap := utils.ConvertProjectConfigToAutoConfigMap(projectConfig)
	s.SetAllConfig(ctx, projectConfigMap)

	//deploy stack
	upRes, err := s.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not deploy stack"})
		logrus.Error("Could not deploy stack", err)
		return
	}

	// Convert output to json and print
	outputJson, err := json.Marshal(upRes.Outputs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not marshal output"})
		return
	}

	fmt.Println(string(outputJson))

	c.JSON(http.StatusCreated, string(outputJson))
	return
	
}

//Deletes a particular stack and all the associated resources
func DeleteStack(c *gin.Context) {
    w.Header().Set("content-type", "application/json")
	
	ctx := context.Background()
	
	stackName := params["stack_name"]
	projectName := params["project_name"]

	s, err := auto.SelectStackInlineSource(ctx, stackName, projectName , PulumiProgram, auto.Program(PulumiProgram))
	if err != nil {
		// check if stack exists
		if auto.IsSelectStack404Error(err){
			w.WriteHeader(404)
			fmt.Fprintf(w, "Stack %v doesn't exist", stackName)
			return
		}
		w.WriteHeader(500)
		fmt.Fprintf(w, fmt.Sprintf(err.Error()))
		return
	}

	// destroy stack
	_, err = s.Destroy(ctx, optdestroy.ProgressStreams(os.Stdout))

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	//delete the stack and all associated history and config
	err = s.Workspace().RemoveStack(ctx, stackName)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(200)
	
}

//GetStack returns a single stack name if it exists
// func GetStack(w http.ResponseWriter, req *http.Request) {

// 	w.Header().Set("content-type", "application/json")

// 	ctx := context.Background()
// 	params := mux.Vars(req)
// 	stackName := params["stack_name"]
// 	projectName := params["project_name"]

// 	_, err := auto.SelectStackInlineSource(ctx, stackName, projectName, auto.Program(createProgram))
// 	if err != nil {
// 		// check if stack exists
// 		if auto.IsSelectStack404Error(err){
// 			w.WriteHeader(404)
// 			fmt.Fprintf(w, "Stack %v doesn't exist", stackName)
// 			return
// 		}
// 		w.WriteHeader(500)
// 		fmt.Fprintf(w, fmt.Sprintf(err.Error()))
// 		return
// 	}

// 	response := &StackResponse{
// 		StackName: stackName,
// 	}

// 	json.NewEncoder(w).Encode(&response)

	
// }
