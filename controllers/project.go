package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Project struct {
	ProjectName string `json:"project_name"`
	AwsRegion   string `json:"aws_region"`
	StackName   string `json:"environment"`
}

type ProjectResponse struct {
	ProjectName string `json:"project_name"`
	AwsRegion   string `json:"aws_region"`
	StackName   string `json:"environment"`
	Status      string `json:"status"`
}

/*
CreateProject: function is responsible for creating a new project on pulumi dashboard.
A Client can have multiple projects.

*/

func CreateProject(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-type", "application/json")

	var project Project
	err := json.NewDecoder(req.Body).Decode(&project)
	if err != nil {
		w.WriteHeader(304)
		fmt.Fprintf(w, "Failed to parse project args")
	}

	err = checkPulumi()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't install Pulumi, ", err), http.StatusInternalServerError)
		return
	}

	pulimiFile, err := os.ReadFile("pulumi-tpl.yaml")
	if err != nil {
		fmt.Println("Could not read file: ", err)
		return
	}

	var pulumiData map[string]interface{}
	err = yaml.Unmarshal(pulimiFile, &pulumiData)
	if err != nil {
		fmt.Println("Could not unmarshal the data: ", err)
		return
	}

	suffixedProjectName := suffixProjectName(project.ProjectName)

	pulumiData["name"] = suffixedProjectName

	//Access the config property
	configProperty, ok := pulumiData["template"].(map[string]interface{})["config"]
	if !ok {
		fmt.Println("Could not find template: config block")
		return
	}

	configProperty.(map[string]interface{})["aws:region"].(map[string]interface{})["default"] = project.AwsRegion
	configProperty.(map[string]interface{})["pulumi:tags"].(map[string]interface{})["projectName"] = suffixedProjectName
	configProperty.(map[string]interface{})["pulumi:tags"].(map[string]interface{})["awsRegionDeployed"] = project.AwsRegion

	pulumiFileBytes, err := yaml.Marshal(pulumiData)
	if err != nil {
		fmt.Println("Could not Marshal the new data: ", err)
	}

	// Create a new pulumi.yaml file in the root directory
	err = os.WriteFile("Pulumi.yaml", pulumiFileBytes, 0644)
	if err != nil {
		fmt.Println("Could not create pulumi.yaml file: ", err)
		return
	}

	_, err = createStackIfDontExist(project.StackName)
	if err != nil {
		fmt.Println("Could not create Stack: ", err)
		http.Error(w, fmt.Sprintf("Error creating Stack %v", err), http.StatusConflict)
	}

	pc_cmd := "pulumi up --stack" + project.StackName + "--skip-preview"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pr, pw := io.Pipe()

	//wait group for go routine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		go func() {
			io.Copy(os.Stdout, pr)			
		}()

		//execute command and write output to pipe
		_, err := ExecuteCommandWithTimeout(ctx, pc_cmd, 30*time.Second, pw)
		pw.CloseWithError(err)

	}()

	response := ProjectResponse {
		Status:      "success",
		AwsRegion:   project.AwsRegion,
		ProjectName: suffixedProjectName,
		StackName:   project.StackName,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}


	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)

}

/*
checkPulumi: Checks if pulumi is installed
*/
func checkPulumi() error {

	//check if pulumi is installed
	checkPulumiCmd := exec.Command("pulumi", "version")
	_, err := checkPulumiCmd.Output()

	if err != nil {
		fmt.Println("Pulumi is not installed")

		//TODO: Maybe install it
		return err
	}
	return nil
}

/*

  ExecuteCommandWithTimeout executes a command with a timeout using Goroutine and context

*/

func ExecuteCommandWithTimeout(ctx context.Context, command string, timeout time.Duration, outputWriter io.Writer) (string, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	cmd.Stderr = outputWriter
	cmd.Stdout = outputWriter
	
	err := cmd.Start()
	if err != nil {
		return "", fmt.Errorf("Command failed to execute ", err)
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
		close(done)
	}()

	select {
	case <- ctx.Done():
		//context cancelled. Terminate command
		cmd.Process.Kill()
		return "", fmt.Errorf("Command terminated")
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("Error Running command ", err)
		}
		return "success", nil
	}

}

/*
A fuction to check if the stack exists, and creates it if doesn't exist
Params:

	stackName: Name of the stack we are checking
*/
func createStackIfDontExist(stackName string) (bool, error) {

	// Check if the stack exists using pulumi stack ls
	cmd := exec.Command("pulumi", "stack", "ls")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error checking stack existence:", err)
		return false, err
	}

	//Check if the stack name already exist
	if !strings.Contains(string(out), stackName) {
		cmd := exec.Command("pulumi", "stack", "init", stackName)
		fmt.Println("Creating Stack...")
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error creating stack:", err)
			return false, err
		}
	}

	return true, nil
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
