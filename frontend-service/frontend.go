package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"okx-bot/frontend-service/app"
	"okx-bot/frontend-service/controllers"
	"okx-bot/frontend-service/models"
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
	start()
}

func start() {
	logger.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts", controllers.GetContactsFor).Methods("GET") //  user/2/contacts

	router.HandleFunc("/api/signal/receive", controllers.ReceiveSignal).Methods("POST") //  user/2/contacts
	router.HandleFunc("/api/signal/create", controllers.CreateSignal).Methods("POST")   //  user/2/contacts
	router.HandleFunc("/api/signal", controllers.GetAllSignals).Methods("GET")          //  user/2/contacts
	router.HandleFunc("/api/signal/bots", controllers.GetBots).Methods("GET")           //  user/2/contacts

	router.HandleFunc("/api/signal/okx/bots", controllers.GetAllOkxBots).Methods("GET")   //  user/2/contacts
	router.HandleFunc("/api/signal/okx/all", controllers.GetAllOkxSignals).Methods("GET") //  user/2/contacts

	router.HandleFunc("/api/bot/create", controllers.CreateBot).Methods("POST")   //  user/2/contacts
	router.HandleFunc("/api/bot/delete", controllers.DeleteBot).Methods("DELETE") //  user/2/contacts

	router.HandleFunc("/api/okx/create", controllers.CreateOkxApi).Methods("POST")
	router.HandleFunc("/api/okx/keys", controllers.GetOkxApiFor).Methods("GET")

	router.HandleFunc("/api/check/okx", controllers.CheckOkx).Methods("GET") //  user/2/contacts

	router.Use(app.JwtAuthentication) //attach JWT auth middleware
	//router.NotFoundHandler = http.NotFoundHandler()
	//GetUserApi

	go func() {
		models.ConnectDB()
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	go func() {
		logger.Infoln("Server REST starting ...")
		err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
		if err != nil {
			logger.Error(err)
		}
		logger.Infoln("Serving REST started")
	}()

	<-c
	logger.Info("Server graceful stopped")
	os.Exit(0)
}
