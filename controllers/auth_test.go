package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"maos-cloud-project-api/models"
	"maos-cloud-project-api/responses"

	// "maos-cloud-project-api/router"
	"maos-cloud-project-api/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type User struct {
	Name     string
	Email    string
	Password string
	Role     string
}

// AuthTestSuite
type AuthTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *AuthTestSuite) SetUpSuite() {

	suite.router = utils.SetUpRouter()
	config := models.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("TEST_DB_NAME"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}

	models.InitDB(config)

}

func (suite *AuthTestSuite) TestSignup() {
	
	log.SetFormatter(&log.TextFormatter{})
	
	correctUser := models.User{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "test1234",
		Role:     "admin",
	}
	
	// Test Case 1: Correct signup
	w := httptest.NewRecorder()

	jsonValue, _ := json.Marshal(correctUser)
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

