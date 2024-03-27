package utils
import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

func SetUpRouter() *gin.Engine{
    router := gin.Default()

	//load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

    return router
}