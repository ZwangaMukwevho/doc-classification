package repository

import "firebase.google.com/go/db"

type FirebaseRepository interface {
	AddWords() error
}

type firebaseClient struct {
	client *db.Client
}

func NewFirebaseRestClient(client *db.Client) firebaseClient {
	return firebaseClient{client}
}

func Get() error {
	return nil
}
