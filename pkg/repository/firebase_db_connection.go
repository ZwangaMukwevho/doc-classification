package repository

import (
	"context"
	"doc-classification/pkg/common"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

func InitDB(databaseURL string, serviceFile string) (*db.Client, error) {
	ctx := context.Background()

	// configure database URL
	conf := &firebase.Config{
		DatabaseURL: databaseURL,
	}

	// fetch service account key
	opt := option.WithCredentialsFile(serviceFile)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		common.Logger.Errorf("error in initializing firebase app: %v", err)
		return nil, err
	}

	client, err := app.Database(ctx)
	if err != nil {
		common.Logger.Errorf("error in creating firebase DB client: %v", err)
		return nil, err
	}

	common.Logger.Info("Successfully initialised DB")
	return client, err
}
