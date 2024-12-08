package resource

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler Handler) *gin.Engine {

	router := gin.Default()

	// Enable CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	router.GET("/ping", handler.pong)
	router.GET("/gmail", handler.initiateGmailAuth)
	router.GET("/gdrive", handler.initiateDriveAuth)
	router.GET("/gmail/authkey", handler.getGmailAuthKey)
	router.POST("/user/create", handler.createUser)
	router.GET("/users", handler.getUsers)
	router.POST("/token/create", handler.createGmailToken)
	router.POST("/token/update", handler.updateToken)

	return router
}
