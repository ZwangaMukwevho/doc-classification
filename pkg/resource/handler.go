package resource

import (
	"context"
	"doc-classification/pkg/common"
	"doc-classification/pkg/model"
	"doc-classification/pkg/repository"
	"doc-classification/pkg/service"
	"log"
	"net/http"
	"os"
	"strings"

	"firebase.google.com/go/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Handler struct {
	// TO-DO: Implement Firebase client in repository
	FirebaseClient      *db.Client
	FirebaseRespository repository.FirebaseRepository
}

// @Summary Ping the api
// @Success 200 string "pong"
// @Failure 500 {string} string "Internal Server Error"
// @Router /ping [get]
func (h *Handler) pong(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "pong")
}

// @Summary Get all words
// @Description Get all words from the Firebase Realtime Database
// @Tags words
// @Produce json
// @Success 200 {array} model.Word
// @Failure 500 {string} string "Internal Server Error"
// @Router /words [get]
func (h *Handler) initiateGmailAuth(c *gin.Context) {
	oAuthByteStream, err := common.GetJsonFileByteStream("google_client_secret.json")
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

	oAuthByteStream, err := common.GetJsonFileByteStream("google_client_secret.json")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "error auth")
		return
	}

	driveConfig, err := google.ConfigFromJSON(*oAuthByteStream, drive.DriveFileScope)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "error auth")
		return
	}
	authURL := service.GetAuthCodeURL(driveConfig)

	c.String(http.StatusOK, authURL)
}

func (h *Handler) getGmailAuthKey(c *gin.Context) {
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

func (h *Handler) createUser(c *gin.Context) {
	var userData model.User

	if err := c.ShouldBindJSON(&userData); err != nil {
		common.Logger.Errorf("Error binding the payload from request: %v", err)
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	gdriveToken, err := service.GetGoogleToken(userData.GdriveCode, drive.DriveReadonlyScope)
	if err != nil {
		common.Logger.Errorf("Error obtaining the Gdrive token on createUser: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	gmailToken, err := service.GetGoogleToken(userData.GmailCode, gmail.GmailReadonlyScope)
	if err != nil {
		common.Logger.Errorf("Error obtaining the Gmail token on createUser: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	driveService, err := initialiseDriveServiceForHandler(gdriveToken)
	if err != nil {
		common.Logger.Errorf("Error initiating the drive service: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	common.Logger.Infof("Initialised the drive service")

	CategoriesInformation := make(map[string]model.Category)
	for _, categoryObject := range userData.Categories {

		directoryId, err := driveService.GetRootDirectoryIDByName(categoryObject.Category)

		if err != nil {
			common.Logger.Errorf("Error obtaining directory id: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		common.Logger.Infof("Direcory ID found: %v", directoryId)

		// If the directory already exists we don't want to create another
		if directoryId != nil {
			CategoriesInformation[*directoryId] = categoryObject
			continue
		}

		folder, err := driveService.CreateDriveDirectory(categoryObject.Category)
		common.Logger.Infof("Found the folder here: %v", folder)

		if err != nil {
			common.Logger.Infof("Error creating the drive directory: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		CategoriesInformation[folder.Id] = categoryObject
	}

	userName := "user_" + strings.ReplaceAll(uuid.NewString(), "-", "")

	var firebaseUser = model.FirebaseUser{
		UserId:     userData.UserId,
		GmailCode:  gmailToken,
		GdriveCode: gdriveToken,
		Categories: CategoriesInformation,
		UserName:   userName,
	}
	common.Logger.Infof("Firebase user initialised: %v", firebaseUser)

	h.FirebaseRespository.UploadUserData(firebaseUser)

	c.IndentedJSON(http.StatusOK, "OK")
}

func (h *Handler) createGmailToken(c *gin.Context) {

	var gmailCode model.Code

	if err := c.ShouldBindJSON(&gmailCode); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	config, err := service.GetOauthConfig(drive.DriveScope)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	tok, err := service.GetTokenUsingAPI(config, gmailCode.CodeString)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, tok)
}

func (h *Handler) getUsers(c *gin.Context) {
	users, err := h.FirebaseRespository.GetUserDataList()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

// getNames godoc
// @Summary      Get all user names
// @Description  Fetches a map of user IDs to user names from Firebase
// @Tags         users
// @Produce      json
// @Success      200  {object}  map[string]string  "Map of userID to userName, e.g., {\"user6ihsvqcxl\": \"Zee430\"}"
// @Failure      500  {object}  map[string]string  "Error message"
// @Router       /users/names [get]
func (h *Handler) getNames(c *gin.Context) {
	userNames := make(map[string]string)

	users, err := h.FirebaseRespository.GetUserDataList()
	common.Logger.Infof("Users are fetched here: %v", users)

	if err != nil {
		common.Logger.Errorf("Error fetching users from firebase repository: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	for userId, user := range *users {
		userNames[userId] = *user.UserName
	}

	common.Logger.Info("Successful fetched users from firebase")
	c.IndentedJSON(http.StatusOK, userNames)
}

func (h *Handler) updateToken(c *gin.Context) {
	var token oauth2.Token

	if err := c.ShouldBindJSON(&token); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	err := h.FirebaseRespository.UpdateGmailToken("mkaax2a72", token)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, "OK")
}

func initialiseDriveServiceForHandler(token *oauth2.Token) (*service.DriveServiceLocal, error) {

	ctx := context.Background()
	b, err := os.ReadFile("google_client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return nil, err
	}

	driveConfig, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}

	driveClient := driveConfig.Client(context.Background(), token)

	driveSrv, err := drive.NewService(ctx, option.WithHTTPClient(driveClient))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
		return nil, err
	}

	srv := service.DriveServiceLocal{Service: driveSrv}
	return &srv, nil
}
