package resource

import (
	"context"
	"doc-classification/pkg/common"
	"doc-classification/pkg/service"
	"fmt"
	"log"
	"net/http"

	"firebase.google.com/go/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/gmail/v1"
)

type Handler struct {
	// TO-DO: Implement Firebase client in repository
	FirebaseClient *db.Client
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
	fmt.Println("auth url: ", authURL)

	c.String(http.StatusOK, authURL)
}

// @Summary Get all words
// @Description Get all words from the Firebase Realtime Database
// @Tags words
// @Produce json
// @Success 200 {array} model.Word
// @Failure 500 {string} string "Internal Server Error"
// @Router /words [get]
func (h *Handler) initiateDriveAuth(c *gin.Context) {

	oAuthByteStream, err := common.GetJsonFileByteStream("client_secret_973692223612-28ae9a7njdsfh7gv89l0fih5q36jt52m.apps.googleusercontent.com.json")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "error auth")
		return
	}

	driveConfig, err := google.ConfigFromJSON(*oAuthByteStream, drive.DriveScope)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "error auth")
		return
	}
	authURL := service.GetAuthCodeURL(driveConfig)
	fmt.Println("auth url: ", authURL)

	c.String(http.StatusOK, authURL)
}

func (h *Handler) getGmailAuthKey(c *gin.Context) {
	var authToken *oauth2.Token

	ref := h.FirebaseClient.NewRef("users/testGmailKey")
	fmt.Println("calling ref")
	if err := ref.Get(context.Background(), &authToken); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	// if !authToken.Valid() { // check if the token is expired
	// 	log.Print("Token not valid")
	// 	c.IndentedJSON(http.StatusInternalServerError, "token is expired")
	// 	return
	// }

	c.IndentedJSON(http.StatusOK, authToken)
}

func (h *Handler) postGmailAuthCode(c *gin.Context) {
	var authToken *oauth2.Token

	ref := h.FirebaseClient.NewRef("users/testGmailKey")
	if err := ref.Get(context.Background(), &authToken); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	// if !authToken.Valid() { // check if the token is expired
	// 	log.Print("Token not valid")
	// 	c.IndentedJSON(http.StatusInternalServerError, "token is expired")
	// 	return
	// }

	c.IndentedJSON(http.StatusOK, authToken)
}
