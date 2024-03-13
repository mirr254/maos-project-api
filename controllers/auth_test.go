package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"maos-cloud-project-api/models"
	"maos-cloud-project-api/responses"
	"maos-cloud-project-api/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type User struct {
	Name     string
	Email    string
	Password string
	Role     string
}

// DatabaseTestSuite
type DatabaseTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *DatabaseTestSuite) SetUpSuite() {

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

func (suite *DatabaseTestSuite) TestSignup() {

	log.SetFormatter(&log.TextFormatter{})

	w := httptest.NewRecorder()
	correctUser := models.User{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "test1234",
		Role:     "admin",
	}

	jsonValue, _ := json.Marshal(correctUser)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusAccepted, w.Code)

	var response responses.UserCreatedResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	suite.Equal("user created", response.Message)

	log.Infof("Signup Response: %s", response.Message )

	noEmail := models.User{
		Name:     "test",
		Password: "test1234",
		Role:     "admin",
	}

	jsonValue, _ = json.Marshal(noEmail)
	req, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusBadRequest, w.Code)

	var noEmailResponse responses.UserCreatedResponse
	json.Unmarshal(w.Body.Bytes(), &noEmailResponse)
	suite.Equal("email must be provided", noEmailResponse.Message)

	loginUser := models.User{
		Name:     "test",
		Password: "test1234",
	}

	jsonValue, _ = json.Marshal(loginUser)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusAccepted, w.Code)

}

func TestDashboard(t *testing.T) {

	mockUnAuthorizedResponse := `{"error":"unauthorized"}`

	r := utils.SetUpRouter()
	r.GET("/dashboard", Dashboard)
	req, _ := http.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockUnAuthorizedResponse, string(responseData))
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func TestServer(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

// func (suite *DatabaseTestSuite) TestLogin() {
// 	testUsers := []models.User{
// 		{Name: "testuser", Email: "test@gmail.com",Password: "test1234", Role: "admin"},
// 	}
// }
