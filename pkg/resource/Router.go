package resource

import "github.com/gin-gonic/gin"

func NewRouter(handler Handler) *gin.Engine {

	router := gin.Default()
	router.GET("/gmail", handler.initiateGmailAuth)
	router.GET("/gdrive", handler.initiateDriveAuth)
	router.GET("/gmail/authkey", handler.getGmailAuthKey)
	router.POST("/gmail/authkey", handler.getGmailAuthKey)

	return router
}
