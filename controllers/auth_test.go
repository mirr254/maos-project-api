package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"maos-cloud-project-api/models"
	// "maos-cloud-project-api/router"
	"maos-cloud-project-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// AuthTestSuite
type AuthTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *AuthTestSuite) SetUpSuite() {

	suite.router = utils.SetUpRouter()
	// router.AuthRoutes(suite.router)
	fmt.Println("Test Suite Setup")

}

func (suite *AuthTestSuite) TestSignup() {
	
	correctUser := models.User{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "test1234",
		Role:     "admin",
	}
	
	// Test Case 1: Correct signup
	w := httptest.NewRecorder()
	jsonValue, _ := json.Marshal(correctUser)

	if suite.router == nil {
		suite.T().Fatal("Test Router not initialized")
	}

	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
	suite.router.ServeHTTP(w, req)
	
	suite.Equal(http.StatusAccepted, w.Code)
	suite.Equal("user created", w.Body.String())

	suite.T().Log("Test Case 1: Correct signup: ",w.Body.String())

	// Test Case 2: correct login
	loginUser := models.User{
		Name:     "test",
		Password: "test1234",
	}

	jsonValue, _ = json.Marshal(loginUser)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusAccepted, w.Code)

	suite.T().Log("Test Case 2: correct login: ",w.Body.String())

	// Test Case 3: check dashboard
	req, _ = http.NewRequest("GET", "/dashboard", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusAccepted, w.Code)
	suite.Equal("admin", w.Body.String())
	
	suite.T().Log("Test Case 3: check dashboard: ",w.Body.String())

	// Test Case 4: no email
	w = httptest.NewRecorder()

	noEmail := models.User{
		Name:     "test",
		Password: "test1234",
		Role:     "admin",
	}

	jsonValue, _ = json.Marshal(noEmail)
	req, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal("email must be provided", w.Body.String())

	suite.T().Log("Test Case 4: no email: ",w.Body.String())

}

func TestServer(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

