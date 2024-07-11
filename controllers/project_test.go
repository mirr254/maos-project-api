package controllers

import (
    "testing"
    // "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    // "github.com/pulumi/pulumi/sdk/v3/go/pulumi/test"
    "github.com/stretchr/testify/suite"
    "strings"
    "strconv"
)

type ProjectTestSuite struct {
    suite.Suite
    project_name string
    aws_region string

}

func (suite *ProjectTestSuite) SetupTest() {
    suite.project_name = "my-project"
    suite.aws_region = "us-west-2"
}

func (suite *ProjectTestSuite) TestSuffixProjectName() {
    suffixedName := suffixProjectName(suite.project_name)

    // Check if the suffixed name starts with the original project name
    if !strings.HasPrefix(suffixedName, suite.project_name) {
        suite.T().Errorf("Expected suffixed name to start with %s, got %s", suite.project_name, suffixedName)
    }

    // Check if the suffixed name has a hyphen followed by a number
    parts := strings.Split(suffixedName, "-")
    if len(parts) != 3 {
        suite.T().Errorf("Expected suffixed name to have two parts separated by a hyphen, got %s", suffixedName)
    }

    _, err := strconv.Atoi(parts[2])
    if err != nil {
        suite.T().Errorf("Expected suffixed name to have a number after the hyphen, got %s", suffixedName)
    }
}

// TestProject is the entry point for this test suite
func TestProject(t *testing.T) {
    suite.Run(t, new(ProjectTestSuite))
}

// func TestPulumiProgram(t *testing.T) {
//     testOptions := &test.ProgramTestOptions{
//         Quick: true,
//         // Sets up the test, and runs the Pulumi program in the context of an isolated environment
//         RunUpdateTest: func(ctx *pulumi.Context, stackInfo test.StackInfo) error {
//             err := PulumiProgram(ctx)
//             assert.NoError(t, err)
//             return err
//         },
//         // Verifies the stack outputs after the update
//         Validate: func(outs pulumi.OutputMap) error {
//             successOutput, ok := outs["Success!"]
//             assert.True(t, ok)
//             assert.Equal(t, "success", successOutput.Value)
//             return nil
//         },
//     }

//     test.RunPulumiProgramTest(t, testOptions)
// }
