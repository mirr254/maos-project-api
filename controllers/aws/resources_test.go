package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CreateVPCTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *CreateVPCTestSuite) SetupTest() {
	suite.router = gin.Default()
	suite.router.POST("/api/v1/aws/vpc", CreateVPCEndpoint)

}

func (suite *CreateVPCTestSuite) prepareTestContext(payload []byte) (*gin.Context, *httptest.ResponseRecorder) {
	
	req, _ := http.NewRequest("POST", "/api/v1/aws/vpc", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	return c, resp
}

func (suite *CreateVPCTestSuite) TestCreateVPC_Success() {
	vpc := map[string]interface{}{
		"project_name": "test-vpc",
		"bastion_host_access_whitelist_cidr": "0.0.0.0/0",
	}
	payload, _ := json.Marshal(vpc)
	suite.T().Log("Payload: ",bytes.NewBuffer(payload))

	c, resp := suite.prepareTestContext(payload)
	CreateVPCEndpoint(c)

	suite.T().Log("Response: ", resp.Body.String())

	assert.Equal(suite.T(), http.StatusCreated, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)

	assert.NotNil(suite.T(), response["vpc"])
}

func TestCreateVPCTestSuite(t *testing.T) {
	suite.Run(t, new(CreateVPCTestSuite))
}