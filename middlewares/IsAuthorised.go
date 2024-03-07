package middlewares

import (
	"maos-cloud-project-api/utils"

	gin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func IsAuthorized() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		cookie, err := ctx.Cookie("token")
		if err != nil {
			ctx.JSON(401, gin.H{"error": "unauthorized"})
			logrus.Error("Error 1: ", err)
			ctx.Abort()
			return
		}

		claims, err := utils.ParseToken(cookie)
		if err != nil {
			ctx.JSON(401, gin.H{"error": "signature is invalid"})
			logrus.Error("Error 2: ", err)
			ctx.Abort()
			return
		}

		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
