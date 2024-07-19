package controllers

// import (
// 	"testing"
// 	"net/http/httptest"

// 	// "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// 	// "github.com/pulumi/pulumi/sdk/v3/go/pulumi/test"

// 	"maos-cloud-project-api/mocks"
// 	"maos-cloud-project-api/utils"

// 	"bytes"
// 	"encoding/json"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/pulumi/pulumi/sdk/v3/go/auto"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// type StackHandlerTestSuite struct {
//     suite.Suite
// 	router           *gin.Engine
// 	w                *httptest.ResponseRecorder
// 	stackHandlerMock *mocks.MockStack

// }

// func (suite *StackHandlerTestSuite) SetupTest() {
//     suite.router = utils.SetUpRouter()
//     suite.router.POST("/api/v1/project/stack", CreateStack)

//     // suite.router.DELETE("api/v1/project/stack", DeleteStack)
//     // suite.project_name = "my-project"
//     // suite.region = "us-west-2"
// }

// func BuildProjectConfig(projectName, region string) map[string]string {
// 	return map[string]string{
// 		"project_name": projectName,
// 		"region":       region,
// 	}
// }

// func ConvertProjectConfigToAutoConfigMap(config map[string]string) auto.ConfigMap {
// 	autoConfig := auto.ConfigMap{}
// 	for k, v := range config {
// 		autoConfig[k] = auto.ConfigValue{Value: v}
// 	}
// 	return autoConfig
// }

// func (suite *StackHandlerTestSuite) Test_CresteStackSuccess() {
// 	mockStack := new(mocks.MockStack)
// 	mockStack.On("SetAllConfig", mock.Anything, mock.Anything).Return(nil)
// 	mockStack.On("Up", mock.Anything, mock.Anything).Return(auto.UpResult{
// 		Outputs: auto.OutputMap{
// 			"output": auto.OutputValue{
// 				Value: "value",
// 			},
// 		},

// 	}, nil)

// 	suite.stackHandlerMock = mockStack

// 	stack := Stack{
// 		StackName:   "test-stack",
// 		ProjectName: "test-project",
// 		Region:      "us-west-2",
// 	}

// 	body, _ := json.Marshal(stack)
// 	req, _ := http.NewRequest("POST", "/api/v1/project/stack", bytes.NewBuffer(body))
// 	resp := httptest.NewRecorder()
// 	suite.router.ServeHTTP(resp, req)

// 	assert.Equal(suite.T(), http.StatusCreated, resp.Code)
// 	mockStack.AssertExpectations(suite.T())
// }

// func TestStack (t *testing.T) {
// 	suite.Run(t, new(StackHandlerTestSuite))
// }
