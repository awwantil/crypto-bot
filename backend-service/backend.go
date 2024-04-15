package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	controllers2 "okx-bot/backend-service/controlers"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.main-rest",
	})
)

func main() {

	logger.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.HandleFunc("/", controllers2.Info).Methods("GET")

	go func() {
		port := "8020"
		err := http.ListenAndServe(":"+port, router)
		if err != nil {
			logger.Error(err)
		}
		logger.Infoln("Backend server started")
	}()

	<-c
	logger.Info("Backend server graceful stopped")
	os.Exit(0)

}
