package repository

import (
	"context"
	"doc-classification/pkg/model"

	"firebase.google.com/go/db"
)

type FirebaseRepository interface {
	UploadUserData(model.FirebaseUser) error
}

type firebaseClient struct {
	client *db.Client
}

func NewFirebaseRestClient(client *db.Client) firebaseClient {
	return firebaseClient{client}
}

func (f firebaseClient) UploadUserData(userdata model.FirebaseUser) error {
	ref := f.client.NewRef("doc-classification/user")

	// Push data to firebase
	if err := ref.Set(context.Background(), userdata); err != nil {
		return err
	}

	return nil
}
