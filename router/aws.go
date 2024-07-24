package router

// import (
// 	"maos-cloud-project-api/controllers"
// 	"github.com/gin-gonic/gin"
// )

// func AWSIamRoutes(r *gin.Engine) {
// 	v1 := r.Group("/api/v1")

	// v1.POST("/:project_name/:stack_name/iam_user", controllers.CreateIAMUser)
	// v1.GET("/:project_name/aws/:stack_name/iam_user", controllers.GetIAMUsers)
	// v1.DELETE("/:project_name/:stack_name/iam_user/<iam_user_arn>", controllers.DeleteIAMUser)
	// v1.PATCH("/:project_name/:stack_name/iam_user/<iam_user_arn>", controllers.UpdateIAMUser)
// }