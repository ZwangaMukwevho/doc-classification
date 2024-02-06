package resource

import (
	"doc-classification/pkg/common"
	"doc-classification/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type Handler struct {
	// TO-DO: Implement Firebase client in repository
	// FirebaseClient *db.Client
}

// @Summary Get all words
// @Description Get all words from the Firebase Realtime Database
// @Tags words
// @Produce json
// @Success 200 {array} model.Word
// @Failure 500 {string} string "Internal Server Error"
// @Router /words [get]
func (h *Handler) initiateGmailAuth(c *gin.Context) {
	oAuthByteStream, err := common.GetJsonFileByteStream("client_secret_973692223612-28ae9a7njdsfh7gv89l0fih5q36jt52m.apps.googleusercontent.com.json")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "error auth")
		return
	}

	gmailConfig, err := google.ConfigFromJSON(*oAuthByteStream, gmail.GmailReadonlyScope)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "error auth")
		return
	}
	authURL := service.GetAuthCodeURL(gmailConfig)

	c.IndentedJSON(http.StatusOK, authURL)
}

// @Summary Get all words
// @Description Get all words from the Firebase Realtime Database
// @Tags words
// @Produce json
// @Success 200 {array} model.Word
// @Failure 500 {string} string "Internal Server Error"
// @Router /words [get]
func (h *Handler) initiateDriveAuth(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, "")
}
