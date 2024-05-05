package controllers // Replace with your actual package name

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	// "maos-cloud-project-api/controllers"
	utils "maos-cloud-project-api/utils"

	"maos-cloud-project-api/mocks"
	"maos-cloud-project-api/models"
)




type SignupTestSuite struct {
	suite.Suite
	router *gin.Engine
	w      *httptest.ResponseRecorder
	c      *gin.Context
	db     *gorm.DB
	
}

func (s *SignupTestSuite) SetupTest() {

	db, err := models.InitDB(config)
	if err != nil {
		// Handle error
		s.T().Fatal("Error initializing database connection")
	}
	db.AutoMigrate(&models.Users{})

    s.router = utils.SetUpRouter()
	s.router.POST("/api/v1/signup", Signup)
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

func (s *SignupTestSuite) Test_InvalidEmail() {

	user := map[string]interface{}{ 
		"name":     "test",
		"email":    "test.gmail.com",
		"password": "plainPassword123", 
		"role":     "admin",  
	  }

	userBody, _ := json.Marshal(user)
	s.T().Log("USER BODY REQ: ", bytes.NewBuffer(userBody))

	ctx, w := s.prepareTestContext(userBody)
	Signup(ctx)

	s.T().Log("RESPONSE BODY: ", w.Body.String())
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Contains(s.T(), w.Body.String(), "invalid email address")
}

func TestSignupSuite(t *testing.T) {
	suite.Run(t, new(SignupTestSuite))
}

func (s *SignupTestSuite) TearDownSuite() {

    result := s.db.Exec("DROP TABLE users")
	if result.Error != nil {
		s.T().Fatal("Error dropping table: ", result.Error)
	} 
	s.T().Log("TestSignupSuite TearDown: dropped table users")

}

type LoginTestSuite struct {
	suite.Suite
	router *gin.Engine
	w      *httptest.ResponseRecorder
	c      *gin.Context
	db     *gorm.DB
}

func (s *LoginTestSuite) SetupTest() {
	
	db, err := models.InitDB(config)
	if err != nil {
		// Handle error
		s.T().Fatal("Error initializing database connection")
	}
	db.AutoMigrate(&models.Users{})

	s.router = utils.SetUpRouter()
	s.router.POST("/api/v1/login", Login)
	s.db = db

	s.w = httptest.NewRecorder()
	s.c, _ = gin.CreateTestContext(s.w)
	
}

func (s *LoginTestSuite) prepareTestContext(userBody []byte) (*gin.Context, *httptest.ResponseRecorder) {
	// Initialize the response recorder
	w := httptest.NewRecorder()

	// Create a new HTTP request with the user body
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(userBody))
	req.Header.Add("Content-Type", "application/json")

	// Create a new gin context from the request
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c, w
}

func (s *LoginTestSuite) Test_ValidLogin() {

	signupUser := map[string]interface{}{ 
		"name":     "test",
		"email":    "test@gmail.com",
		"password": "plainPassword123", 
		"role":     "admin", 
	  }
	userBody, _ := json.Marshal(signupUser)

	ctx, w := s.prepareTestContext(userBody)
	Signup(ctx)
	s.T().Log("LOGIN_SUITE: Signup RESPONSE BODY: ", w.Body.String())

	loginUser := map[string]interface{}{	
		"email":    "test@gmail.com",
		"password": "plainPassword123",
	}
	loginBody, _ := json.Marshal(loginUser)
	s.T().Log("USER BODY REQ: ", bytes.NewBuffer(loginBody))

	ctx, w = s.prepareTestContext(loginBody)
	Login(ctx)

	s.T().Log("Login RESPONSE BODY: ", w.Body.String())
	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Contains(s.T(), w.Body.String(), "user logged in")
}

func (s *LoginTestSuite) Test_InvalidEmailLogin() {

	signupUser := map[string]interface{}{ 
		"name":     "test",
		"email":    "test@gmail.com",
		"password": "plainPassword123", 
		"role":     "admin", 
	  }
	userBody, _ := json.Marshal(signupUser)
	s.T().Log("USER BODY REQ: ", bytes.NewBuffer(userBody))

	ctx, w := s.prepareTestContext(userBody)
	Signup(ctx)
	s.T().Log("Signup RESPONSE BODY: ", w.Body.String())

	LoginUser := map[string]interface{}{ 
		"email":    "email@me",
		"password": "plainPassword123",
	}
	loginBody, _ := json.Marshal(LoginUser)
	s.T().Log("USER BODY REQ: ", bytes.NewBuffer(loginBody))

	ctx, w = s.prepareTestContext(loginBody)
	Login(ctx)

	s.T().Log("RESPONSE BODY: ", w.Body.String())
	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(s.T(), w.Body.String(), "invalid username or password")

}

func (s *LoginTestSuite) Test_InvalidPasswordLogin() {

	signupUser := map[string]interface{}{ 
		"name":     "test",
		"email":    "test@gmail.com",
		"password": "plainPassword123", 
		"role":     "admin", 
	  }

	userBody, _ := json.Marshal(signupUser)
	s.T().Log("USER BODY REQ: ", bytes.NewBuffer(userBody))

	ctx, w := s.prepareTestContext(userBody)
	Signup(ctx)
	s.T().Log("Signup RESPONSE BODY: ", w.Body.String())
	
	LoginUser := map[string]interface{}{
		"email":    "test@gmail.com",
		"password": "plainPassw",
	}
	loginBody, _ := json.Marshal(LoginUser)
	s.T().Log("USER BODY REQ: ", bytes.NewBuffer(loginBody))
		
	ctx, w = s.prepareTestContext(loginBody)
	Login(ctx)

	s.T().Log("RESPONSE BODY: ", w.Body.String())
	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(s.T(), w.Body.String(), "invalid username or password")

}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (s *LoginTestSuite) TearDownSuite() {

	s.db.Exec("DROP TABLE users")
	s.T().Log("TearDownSuite")

}

type EmailVerficationLinkTestSuite struct {
	suite.Suite
	router          *gin.Engine
	db 			    *gorm.DB
	w               *httptest.ResponseRecorder
	c               *gin.Context
	mockEmailSender *mocks.MockEmailSender
}

func (s *EmailVerficationLinkTestSuite) prepareTestContext(userBody []byte) (*gin.Context, *httptest.ResponseRecorder) {
	// Initialize the response recorder
	w := httptest.NewRecorder()

	// Create a new HTTP request with the user body
	req := httptest.NewRequest("POST", "/api/v1/send-verification-email", bytes.NewBuffer(userBody))
	req.Header.Add("Content-Type", "application/json")

	// Create a new gin context from the request
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c, w
}

func (s *EmailVerficationLinkTestSuite) SetupTest() {

	db, err := models.InitDB(config)
	if err != nil {
		// Handle error
		s.T().Fatal("Error initializing database connection")
	}
	db.AutoMigrate(&models.Users{})
	s.mockEmailSender = new(mocks.MockEmailSender)

	s.router = utils.SetUpRouter()
	s.router.POST("/api/v1/send-verification-email", func(c *gin.Context) {
		SendEmailVerificationLink(s.mockEmailSender)(c)
	})
	s.db = db

	s.w = httptest.NewRecorder()
	s.c, _ = gin.CreateTestContext(s.w)

}

func (suite *EmailVerficationLinkTestSuite) Test_SendEmailVerificationLinkSuccess() {

	signupUser := map[string]interface{}{ 
		"name":     "test",
		"email":    "test@gmail.com",
		"password": "plainPassword123", 
		"role":     "admin", 
	  }
	userBody, _ := json.Marshal(signupUser)

	ctx, w := suite.prepareTestContext(userBody)
	Signup(ctx)
	suite.T().Log("EMAIL_VERIFICATION_LINK_SUITE: Signup RESPONSE BODY: ", w.Body.String())

	loginUser := map[string]interface{}{	
		"email":    "test@gmail.com",
		"password": "plainPassword123",
	}
	loginBody, _ := json.Marshal(loginUser)

	ctx, w = suite.prepareTestContext(loginBody)
	Login(ctx)

	var token string
	cookies := w.Result().Cookies()

	for _, cookie := range cookies {
		if cookie.Name == "token" {
			token = cookie.Value
			break
		}
	}

	if token == "" {
		suite.T().Fatal("Token not found in cookies")
	}
		
	// email_verification_token, err := utils.GenerateToken()
	// if err != nil {
	// 	suite.T().Fatal("Failed to generate token:", err)
	// }
	// subject := "Email Verification"
	// route := "verify-email"
	// body := "Click the link below to verify your email\n" + utils.CreateVerificationLink(route, email_verification_token)
	suite.mockEmailSender.On("SendEmail",
		mock.MatchedBy(func(email string) bool { return email == "test@gmail.com" }),
		mock.MatchedBy(func(subject string) bool { return subject == "Email Verification" }),
		mock.MatchedBy(func(body string) bool { return strings.Contains(body, "Click the link below to verify your email") }),
	).Return(nil)

	emailPayload := map[string]string{
		"email": "test@gmail.com",
	}
	emailBody, err := json.Marshal(emailPayload)
	if err != nil {
		suite.T().Fatal("Failed to marshal JSON:", err)
	}
	
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/send-verification-email", bytes.NewBuffer(emailBody))
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.T().Log("EMAIL_VERIFICATION_LINK_SUITE RESPONSE BODY: ", resp.Body.String())

	suite.Equal(http.StatusOK, resp.Code)
	// suite.mockEmailSender.AssertExpectations(suite.T())

}

// func (suite *EmailVerficationLinkTestSuite) Test_SendEmailVerificationAlreadyVerified(){
// 	suite.user.IsEmailVerified = true
// 	req, _ := http.NewRequest(http.MethodPost, "/api/v1/send-verification-email", nil)
// 	resp := httptest.NewRecorder()
// 	suite.router.ServeHTTP(resp, req)

// 	suite.Equal(http.StatusBadRequest, resp.Code)
// }

func TestEmailVerficationLinkTestSuite(t *testing.T) {
	suite.Run(t, new(EmailVerficationLinkTestSuite))
}

func (s *EmailVerficationLinkTestSuite) TearDownSuite() {

	// s.db.Exec("DROP TABLE users")
	s.T().Log("TearDownSuite")

}