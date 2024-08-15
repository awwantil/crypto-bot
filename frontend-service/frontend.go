package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"okx-bot/exchange/model"
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

		logger.Infoln("Execute GetOrderInfo")
		opts := model.OptionParameter{
			Key:   "contractAlias",
			Value: "SWAP",
		}
		api, err := app.GetOkxApi(1)
		if err != nil {
			logger.Errorf("Error in GetOkxApi: %v", err)
		}
		//
		//id := "1695898177540702208"
		//req := new(model.BaseOrderRequest)
		//req.InstId = "SOL-USDT-SWAP"
		//req.OrdId = id
		//resp3, resp4, err := api.Isolated.GetOrderInfo(*req, opts)
		//logger.Info("err: ", err)
		//logger.Info("resp3: ", resp3)
		//logger.Info("resp4: ", string(resp4))

		req := new(model.SetLeverageRequest)
		req.InstId = "SOL-USDT-SWAP"
		req.Lever = "3.02"
		req.MgnMode = model.ISOLATED
		//req.Ccy = "USDT"
		//req.PosSide = model.LONG
		resp3, resp4, err := api.Isolated.SetLeverage(*req, opts)
		logger.Info("err: ", err)
		logger.Info("resp3: ", resp3)
		logger.Info("resp4: ", string(resp4))

		orderId, err := app.CreateOrder(api, "SOL-USDT", "5")
		if err != nil {
			return
		}
		logger.Infoln(orderId)

	}()

	go func() {
		//app.InitOkxApi()
		//app.StartDeal("SOL-USDT", "short")
		//resp1, resp2, err := app.GetApi().Isolated.GetAccount("SOL")
		//logger.Info("err: ", err)
		//logger.Info("resp1: ", resp1)
		//logger.Info("resp2: ", string(resp2))
		//1345253813988876288

		//posHistoryRequest := new(model.FuturesPositionHistoryRequest)
		//posHistoryRequest.InstId = "SOL-USDT-SWAP"
		//posHistory, _, err := app.GetApi().GetPositionsHistory(*posHistoryRequest)
		//if err != nil {
		//	panic(err)
		//}
		//logger.Info("posHistory = ", posHistory)
		//logger.Info("posHistory = ", posHistory[0].Pnl)
		//logger.Info("posHistory = ", posHistory[0].RealizedPnl)
		//logger.Info("posHistory = ", posHistory[0].Type)

		//app.EndDeal("SOL-USDT", "1712585104400", "")
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
