package project

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// create a function that accepts project name, pulumi.yaml file for project defination
// exec pulumi initial commands using go os libraries


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

	//Changing the default project values
    pulumiData["name"] = suffixProjectName(projectName)

	pulumiData["template"].(map[string]interface{})["config"].(map[string]interface{}).(map[string]interface{})["aws:region"]["default"].(map[string]interface{})["region"] = awsregion

	pulumiFileBytes, err := yaml.Marshal(pulumiData)
	if err != nil {
		fmt.Println("Could not Marshal the new data: ", err)
	}

	//TODO: Remove
	fmt.Println(string(pulumiFileBytes))
	fmt.Println("Region: ", awsregion)

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