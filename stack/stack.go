package stack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Stack struct {
	StackName  string `json:"stack_name"`
	ProjectName string `json:"project_name"`
}

type StackResponse struct {
	StackName  string `json:"stack_name"`
}

// This package handles all stack operations including, deleting, updating, list/get stack

// Creates a new stack/project
// Return the stack name if created successfully
func CreateStack(w http.ResponseWriter, req *http.Request) {

    w.Header().Set("Content-type", "application/json")
    
	var stack Stack
    err := json.NewDecoder(req.Body).Decode(&stack)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Failed to parse stack args")
		return
	}

	// for debugging purposes
	stackData, _ := json.MarshalIndent(stack, "", "\t")
	fmt.Println(string(stackData))
    
	ctx := context.Background()


	s, err := auto.NewStackInlineSource( ctx, stack.StackName, stack.ProjectName, createProgram() )
	
	if err != nil {
		//check if stack already exist
		if auto.IsCreateStack409Error(err) {
			w.WriteHeader(409)
			fmt.Fprintf(w, fmt.Sprintf( "Stack %v in project %v already exists.", stack.StackName, stack.ProjectName ))
			return
		}

		w.WriteHeader(500)
		fmt.Fprintf(w, "An error %v occurred",err.Error())
		return
	}

	s.SetConfig(ctx, "project_name", auto.ConfigValue{Value: stack.ProjectName})
	s.SetConfig(ctx, "stack_name", auto.ConfigValue{ Value: stack.StackName })

	response := &Stack{
		StackName: stack.StackName,
		ProjectName: stack.ProjectName,
	}
	json.NewEncoder(w).Encode(&response)
	
}

//Deletes a particular stack and all the associated resources
func DeleteStack(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("content-type", "application/json")
	
	ctx := context.Background()
	params := mux.Vars(req)
	stackName := params["stack_name"]
	projectName := params["project_name"]

	s, err := auto.SelectStackInlineSource(ctx, stackName, projectName , createProgram())
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
func GetStack(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("content-type", "application/json")

	ctx := context.Background()
	params := mux.Vars(req)
	stackName := params["stack_name"]
	projectName := params["project_name"]

	_, err := auto.SelectStackInlineSource(ctx, stackName, projectName, createProgram())
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

	response := &StackResponse{
		StackName: stackName,
	}

	json.NewEncoder(w).Encode(&response)

	
}

func createProgram() pulumi.RunFunc{
	return func(ctx *pulumi.Context) error {

		return nil 
	}
}