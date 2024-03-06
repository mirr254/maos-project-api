package middlewares

import (
	"maos-cloud-project-api/utils"

	gin "github.com/gin-gonic/gin"
)

func IsAuthorized() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		cookie, err := ctx.Cookie("token")
		if err != nil {
			ctx.JSON(401, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		claims, err := utils.ParseToken(cookie)
		if err != nil {
			ctx.JSON(401, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
