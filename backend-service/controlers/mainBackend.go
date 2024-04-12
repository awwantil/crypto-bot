package controllers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	u "okx-bot/frontend-service/utils"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.main-backend-controllers",
	})
)

var Info = func(w http.ResponseWriter, r *http.Request) {
	logger.Infoln("Get info")
	u.Respond(w, u.Message(true, "Service is working now"))
}
