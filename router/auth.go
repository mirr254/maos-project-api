package router

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	controllers "maos-cloud-project-api/controllers"
	"maos-cloud-project-api/middlewares"
	models "maos-cloud-project-api/models"
	utils "maos-cloud-project-api/utils"

	"github.com/gin-gonic/gin"
)

var jwtkey = []byte(os.Getenv("SECRET_KEY"))

func AuthRoutes(r *gin.Engine) {
	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
	r.GET("/dashboard", middlewares.IsAuthorized(), controllers.Dashboard)
	r.GET("/logout",middlewares.IsAuthorized(), controllers.Logout)
}

func Login(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	var existingUser models.User
	models.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{"error": "User doesn't exist"})
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)
	if !errHash {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}

	expiration_time := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expiration_time.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString(jwtkey)

	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}

	c.SetCookie("token", token_string, int(expiration_time.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged in"})
}
