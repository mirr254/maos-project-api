package controllers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	models "maos-cloud-project-api/models"
	utils "maos-cloud-project-api/utils"

	gin "github.com/gin-gonic/gin"
)

type ErrorResponse struct {
    Error string `json:"error"`
}

func Signup(c *gin.Context) {

    var user models.User
	config := utils.GetEnvVars()
	db, err := models.InitDB(config)
	if err != nil {
		// Handle error
		panic(err)
	}

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
        return
    }

    // Check for empty email
    if user.Email == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "email must be provided"})
        return
    }

    var existingUser models.User


    db.Where("email = ?", user.Email).First(&existingUser)

    if existingUser.ID != 0 {
        c.JSON(http.StatusConflict, ErrorResponse{Error: "user already exists"})
        return
    }

    hashedPassword, err := utils.GenerateHashPassword(user.Password)

    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "could not generate password hash"})
        return
    }

    user.Password = hashedPassword

    db.Create(&user)

    c.JSON(http.StatusCreated, gin.H{"success": "user created"})
	
	return 
}

func Login(c *gin.Context) {

	var user models.User
	config := utils.GetEnvVars()
	db, err := models.InitDB(config)
	if err != nil {
		// Handle error
		panic(err)
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	db.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey, err := utils.GetSecretKey()
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	jwtKey := []byte(secretKey)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"success": "user logged in"})
}

func ResetPassword(c *gin.Context) {

	var user models.User

	config := utils.GetEnvVars()
	db, err := models.InitDB(config)
	if err != nil {
		// Handle error
		panic(err)
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	db.Where("email = ?", user.Email).First(&existingUser)
	// models.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		c.JSON(500, gin.H{"error": "could not generate password hash"})
		return
	}

	db.Model(&existingUser).Update("password", user.Password)

	c.JSON(http.StatusOK, gin.H{"success": "password updated"})
}

func Dashboard(c *gin.Context) {

	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, err := utils.ParseToken(cookie)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if claims.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "customer dashboard", "role": claims.Role})

}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"success": "user logged out"})
}
