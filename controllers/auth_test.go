package controllers // Replace with your actual package name

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
	
	utils "maos-cloud-project-api/utils"

	models "maos-cloud-project-api/models"
)

type MockHarsher struct {
    mock.Mock
}


type SignupTestSuite struct {
	suite.Suite
	router *gin.Engine
	w      *httptest.ResponseRecorder
	c      *gin.Context
	
}

func (s *SignupTestSuite) SetupTest() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Env file not loaded. Exiting...", err )
	} 

	r := utils.SetUpRouter()

	s.router = r // gin.Default()
	gin.SetMode(gin.TestMode)
	s.w = httptest.NewRecorder()
	s.router.POST("/signup", Signup)

	s.c, _ = gin.CreateTestContext(s.w)
	

}

func (m *MockHarsher) GenerateHashPassword(password string) (string, error) {
    args := m.Called(password)
    return args.String(0), args.Error(1)
}

// SignupWithMockHasher is a wrapper function for testing Signup with a mock hasher
func SignupWithMockHasher(c *gin.Context, hasher MockHarsher, user models.User) (string, error) {
	hashedPassword, err := hasher.GenerateHashPassword(user.Password)
	if err != nil {
	  return "", err
	}

	user.Password = hashedPassword
	Signup(c)

	return hashedPassword, nil
  }

func (s *SignupTestSuite) Test_ValidSignup() {
	//use the hasher wraper function
	var mockHarsher MockHarsher
	mockHarsher.On("GenerateHashPassword", "test1234" ).Return("hashedPassword", nil)
	hashedPassword, err := SignupWithMockHasher(s.c, mockHarsher, models.User{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "test1234",
		Role:     "admin",
	})

	user := map[string]string{ // Use string keys for field names
		"name":     "test",
		"email":    "test@gmail.com",
		"password": "plainPassword123", // Replace with actual password
		"role":     "admin", // Optional
	  }
	userBody, _ := json.Marshal(user)
	
	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(userBody))

	s.T().Log("USER BODY REQ: ", bytes.NewBuffer(userBody))
	

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "hashedPassword", hashedPassword)
	
	s.T().Log("REQUEST : ",req)

	s.router.ServeHTTP(s.w, req)
	
	assert.Equal(s.T(), 201, s.w.Code)

	// var actualBody map[string]interface{}
	// err = json.Unmarshal(s.w.Body.String(), &actualBody)
	s.T().Log("RESPONSE BODY: ", s.w.Body.String())
	assert.NoError(s.T(), err)
	// assert.Equal(s.T(), map[string]interface{}{"success": "user created"}, actualBody)
}

func (s *SignupTestSuite) Test_EmptyEmail() {
	requestBody := `{"password": "password123"}`
	req := httptest.NewRequest("POST", "/signup", bytes.NewReader([]byte(requestBody)))

	s.router.POST("/signup", Signup)
	s.router.ServeHTTP(s.w, req)

	assert.Equal(s.T(), 400, s.w.Code)
	// assert.Equal(s.T(), map[string]interface{}{" error ": ` email must be provided `}, s.w.Body.String())
}

// Add similar test cases for Existing Email and Invalid JSON as before

func TestSignupSuite(t *testing.T) {
	suite.Run(t, new(SignupTestSuite))
}

func (s *SignupTestSuite) TearDownSuite() {
    // Close database connection and clean up resources
    sqlDB, err := models.DB.DB()
	if err != nil {
		s.T().Fatal("Error closing the database connection")
	}
	s.T().Log("Closing the database connection")
	defer sqlDB.Close()
}
