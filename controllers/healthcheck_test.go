package controllers

import (
	"testing"
	"net/http"
	"net/http/httptest"
	
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"maos-cloud-project-api/utils"
	
)

type HealthCheckTestSuite struct {
	suite.Suite
	router *gin.Engine
	w 	*httptest.ResponseRecorder
	c 	*gin.Context
}

func (suite *HealthCheckTestSuite) SetupTest() {

	suite.router = utils.SetUpRouter()
	suite.router.GET("/api/v1/health", HealthCheck)
	suite.w = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	suite.c, _ = gin.CreateTestContext(suite.w)
	suite.c.Request = req
}

func (suite *HealthCheckTestSuite) TestHealthCheckCode() {
	suite.router.ServeHTTP(suite.w, suite.c.Request)
	suite.Equal(http.StatusOK, suite.w.Code)
}

func (suite *HealthCheckTestSuite) TestHealthCheckBody() {
	suite.router.ServeHTTP(suite.w, suite.c.Request)
	suite.Equal("{\"message\":\"HEALTHY\"}", suite.w.Body.String())
}


// TestHealthCheck is the entry point for this test suite
func TestHealthCheck(t *testing.T) {
	suite.Run(t, new(HealthCheckTestSuite))
}