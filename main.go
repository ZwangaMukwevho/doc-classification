package main

import (
	"doc-classification/pkg/common"
	cronJob "doc-classification/pkg/cron"
	"doc-classification/pkg/repository"
	"doc-classification/pkg/resource"
	"log"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func main() {

	// Load environment variables
	// Load environment variables from .env file
	go setupCron()
	common.InitLogger()
	common.Logger.Info("Initialised logger")

	firebaseDB, err := repository.InitDB("https://react-getting-started-78f85-default-rtdb.firebaseio.com", "firebase_service.json")
	if err != nil {
		log.Fatal(err)
	}

	basePath := ":8000"

	firebaseRepository := repository.NewFirebaseRestClient(firebaseDB)

	router := resource.NewRouter(
		resource.Handler{
			FirebaseClient:      firebaseDB,
			FirebaseRespository: firebaseRepository,
		},
	)

	router.Run(basePath)
}

func setupCron() {
	c := cron.New()

	// Schedule the job to run every minute
	// */3 * * * * fixing
	// 0 0 * * * normal
	_, err := c.AddFunc("0 0 * * *", cronJob.ClassificationCron)
	if err != nil {
		log.Println("Error scheduling cron job:", err)
		return
	}

	c.Start()

	// Run the cron scheduler in the background
	select {}
}
