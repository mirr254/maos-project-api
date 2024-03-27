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

	mocks "maos-cloud-project-api/mocks"
	models "maos-cloud-project-api/models"
)

type SignupTestSuite struct {
	suite.Suite
	router *gin.Engine
	c     *gin.Context
	w     *httptest.ResponseRecorder
	
}

func (s *SignupTestSuite) SetupTest() {
	s.router = gin.New()
	gin.SetMode(gin.TestMode)
	s.w = httptest.NewRecorder()
}

func (s *SignupTestSuite) Test_ValidSignup() {
	correctUser := models.User{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "test1234",
		Role:     "admin",
	}
	jsonValue, _ := json.Marshal(correctUser)
	
	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	
	var mockHarsher mocks.MockHarsher
	mockHarsher.On("GenerateHashPassword", mock.Anything ).Return("hashedPassword", nil)

	// TODO: Solve the error in the next line
	//use the hasher wraper function
	hashedPassword, err := SignupWithMockHasher(s.c, mockHarsher, models.User{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "test1234",
		Role:     "admin",
	})

	s.T().Log("TEST: ",hashedPassword)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "hashedPassword", hashedPassword)
	
	s.router.ServeHTTP(s.w, req)

	assert.Equal(s.T(), 201, s.w.Code)

	var actualBody map[string]interface{}
	err = json.Unmarshal(s.w.Body.Bytes(), &actualBody)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), map[string]interface{}{"success": "user created"}, actualBody)
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
