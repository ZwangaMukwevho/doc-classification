package resource

import "github.com/gin-gonic/gin"

func NewRouter(handler Handler) *gin.Engine {

	router := gin.Default()
	router.GET("/gmail", handler.initiateGmailAuth)
	router.GET("/drive", handler.initiateDriveAuth)

	return router
}
