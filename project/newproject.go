package project

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// A function to run pulumi cli commands

func PulumiCli() {

}
/*
CreateProject: function is responsible for creating a new project on pulumi dashboard.
A Client can have multiple projects.
params: projectName: This can be client's new project

*/
func CreateProject( projectName string, region *string ) ([]byte, error) {

	var awsregion string
	if region != nil {
        awsregion = *region
	}
	
	pulimiFile, err := os.ReadFile("pulumi.yaml")
	if err != nil {
		fmt.Println("Could not read file: ", err)
		return nil, err
	}

	var pulumiData map[string]interface{}
	err = yaml.Unmarshal(pulimiFile, &pulumiData)
	if err != nil {
       fmt.Println("Could not unmarshal the data: ", err)
	   return nil, err
	}

    pulumiData["name"] = suffixProjectName(projectName)
    
	//Access the config property 
	configProperty, ok := pulumiData["template"].(map[string]interface{})["config"]
	if !ok {
       fmt.Println("Could not find aws:region block")
	   return nil, fmt.Errorf("Could not find aws:region block")
	}

	configProperty.(map[string]interface{})["aws:region"].(map[string]interface{})["default"] = awsregion
	configProperty.(map[string]interface{})["pulumi:tags"].(map[string]interface{})["projectName"] = suffixProjectName(projectName)
	configProperty.(map[string]interface{})["pulumi:tags"].(map[string]interface{})["awsRegionDeployed"] = suffixProjectName(projectName)


	pulumiFileBytes, err := yaml.Marshal(pulumiData)
	if err != nil {
		fmt.Println("Could not Marshal the new data: ", err)
	}

	//Debugging purposes
	fmt.Println("Region: ", awsregion)
	fmt.Println(string(pulumiFileBytes))

	return pulumiFileBytes, nil

}

/*

This function adds a suffix to the provided project name to avoid duplicate names
Params: projectName
Returns: suffix-projectname

*/
func suffixProjectName( projectName string) string {

	rand.Seed( time.Now().UnixNano() )
	min := 100
	max := 10000

	return projectName + "-" + strconv.Itoa(rand.Intn( max - min + 1 ) ) 
}