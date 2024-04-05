package controllers // Replace with your actual package name

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	// "maos-cloud-project-api/controllers"
	utils "maos-cloud-project-api/utils"

	models "maos-cloud-project-api/models"
)


type SignupTestSuite struct {
	suite.Suite
	router *gin.Engine
	w      *httptest.ResponseRecorder
	c      *gin.Context
	db     *gorm.DB
	
}

func (s *SignupTestSuite) SetupTest() {

	err := godotenv.Load()
    if err != nil {
        logrus.Fatal("Env file not loaded. Exiting...", err)
    }

	config := utils.GetEnvVars()
	db, err := models.InitDB(config)
	if err != nil {
		// Handle error
		s.T().Fatal("Error initializing database connection")
	}
	db.AutoMigrate(&models.User{})

    s.router = utils.SetUpRouter()
	s.router.POST("/signup", Signup)
	s.db = db

	s.w = httptest.NewRecorder()
	s.c, _ = gin.CreateTestContext(s.w)
	

}

// Simulate a HTTP request with a user body
func (s *SignupTestSuite) prepareTestContext(userBody []byte) (*gin.Context, *httptest.ResponseRecorder) {
    // Initialize the response recorder
    w := httptest.NewRecorder()

    // Create a new HTTP request with the user body
    req := httptest.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(userBody))
    req.Header.Add("Content-Type", "application/json")

    // Create a new gin context from the request
    c, _ := gin.CreateTestContext(w)
    c.Request = req

    return c, w
}


func (s *SignupTestSuite) Test_ValidSignup() {

	user := map[string]interface{}{ 
		"name":     "test",
		"email":    "test@gmail.com",
		"password": "plainPassword123", 
		"role":     "admin", 
	  }
	userBody, _ := json.Marshal(user)
	s.T().Log("USER BODY REQ: ", bytes.NewBuffer(userBody))

	ctx, w := s.prepareTestContext(userBody)
	Signup(ctx)

	s.T().Log("RESPONSE BODY: ", w.Body.String())
	assert.Equal(s.T(), http.StatusCreated, w.Code)
	assert.Contains(s.T(), w.Body.String(), "user created")

}

func (s *SignupTestSuite) Test_EmptyEmail() {

	user := map[string]interface{}{ 
		"name":     "test",
		"password": "plainPassword123", 
		"role":     "admin", 
	  }

	userBody, _ := json.Marshal(user)
	s.T().Log("USER BODY REQ: ", bytes.NewBuffer(userBody))

	ctx, w := s.prepareTestContext(userBody)
	Signup(ctx)

	s.T().Log("RESPONSE BODY: ", w.Body.String())
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Contains(s.T(), w.Body.String(), "email must be provided")
}

func TestSignupSuite(t *testing.T) {
	suite.Run(t, new(SignupTestSuite))
}

func (s *SignupTestSuite) TearDownSuite() {

    s.db.Exec("DROP TABLE users")
	s.T().Log("TearDownSuite")

}
