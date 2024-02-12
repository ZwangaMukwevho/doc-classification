package repository

import (
	"context"
	"log"

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
		log.Fatalln("error in initializing firebase app: ", err)
		return nil, err
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase DB client: ", err)
		return nil, err
	}

	log.Print("Successfully initialised ")
	return *&client, err
}
