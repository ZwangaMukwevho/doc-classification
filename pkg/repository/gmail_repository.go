package repository

import (
	"context"
	"doc-classification/pkg/model"
	"log"

	"firebase.google.com/go/db"
	"golang.org/x/oauth2"
)

type FirebaseRepository interface {
	UploadUserData(model.FirebaseUser) error
	GetUserDataList() (*map[string]model.Users, error)
	UpdateGmailToken(userId string, token oauth2.Token) error
	UpdateGdriveToken(userId string, token oauth2.Token) error
}

type firebaseClient struct {
	client *db.Client
}

func NewFirebaseRestClient(client *db.Client) firebaseClient {
	return firebaseClient{client}
}

func (f firebaseClient) UploadUserData(userdata model.FirebaseUser) error {
	ref := f.client.NewRef("doc-classification/user" + userdata.UserId)

	// Push data to firebase
	if err := ref.Set(context.Background(), userdata); err != nil {
		return err
	}

	return nil
}

func (f firebaseClient) GetUserDataList() (*map[string]model.Users, error) {
	ref := f.client.NewRef("doc-classification")

	// Read data from firebase
	var users map[string]model.Users
	if err := ref.Get(context.Background(), &users); err != nil {
		log.Printf("Error reading from database: %v", err)
		return nil, err
	}

	return &users, nil
}

func (f firebaseClient) UpdateGmailToken(userId string, token oauth2.Token) error {
	ref := f.client.NewRef("doc-classification/user" + userId + "/GmailCode")

	// Push data to firebase
	if err := ref.Set(context.Background(), token); err != nil {
		return err
	}

	return nil
}

func (f firebaseClient) UpdateGdriveToken(userId string, token oauth2.Token) error {
	ref := f.client.NewRef("doc-classification/user" + userId + "/GdriveCode")

	// Push data to firebase
	if err := ref.Set(context.Background(), token); err != nil {
		return err
	}

	return nil
}
