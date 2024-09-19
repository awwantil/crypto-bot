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
		//opts := model.OptionParameter{
		//	Key:   "contractAlias",
		//	Value: "SWAP",
		//}
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

		//req := new(model.SetLeverageRequest)
		//req.InstId = "SOL-USDT-SWAP"
		//req.Lever = "3.02"
		//req.MgnMode = model.ISOLATED
		////req.Ccy = "USDT"
		////req.PosSide = model.LONG
		//resp3, resp4, err := api.Isolated.SetLeverage(*req, opts)
		//logger.Info("err: ", err)
		//logger.Info("resp3: ", resp3)
		//logger.Info("resp4: ", string(resp4))

		//newSignal, err := app.CreateSignal(api, "SOL_15m_TestStrategy", "Test strategy")
		//if err != nil {
		//	return
		//}
		//logger.Infoln(newSignal.SignalChanToken)
		//data[{"signalChanId":"1799976084486750208","signalChanToken":"ThclzGQB2McTgH3bPwORvcZcB0aU/KVJ5dbcI7OjfGvGD9jpd46aUBJLt2ZwVvMlQQGIpKcCakUdgMaKlPKqtg=="}]

		//newSignalBot, err := app.CreateSignalBot(api, "1799976084486750208", "SOL-USDT-SWAP", "3", "60")
		//if err != nil {
		//	return
		//}
		////data[{"algoClOrdId":"","algoId":"1816982685798109184","sCode":"0","sMsg":""}]
		////responseBody{"code":"0","data":[{"algoClOrdId":"","algoId":"1817009927903252480","sCode":"0","sMsg":""}],"msg":""}
		//logger.Infof("newSignalBot.AlgoId = %s", newSignalBot.AlgoId)
		//1816982685798109184
		//1817009927903252480
		//1817229891331424256
		//1817234945836847104
		//1817240302030159872
		//1817256141265571840
		//1817262238944722944
		//1817274404305371136

		//_, err = app.CancelSignalBot(api, "1811814804349255680")
		////{"code":"0","data":[{"algoClOrdId":"","algoId":"1811814804349255680","sCode":"0","sMsg":""}],"msg":""}
		//if err != nil {
		//	return
		//}

		placeSubOrderSignalBotRequest := new(model.PlaceSubOrderSignalBotRequest)
		placeSubOrderSignalBotRequest.InstId = "SOL-USDT-SWAP"
		placeSubOrderSignalBotRequest.AlgoId = "1817262238944722944"
		placeSubOrderSignalBotRequest.Side = "buy"
		placeSubOrderSignalBotRequest.OrdType = "market"
		placeSubOrderSignalBotRequest.Sz = "2.00"

		//_, err = app.PlaceSubOrderSignalBot(api, placeSubOrderSignalBotRequest)
		if err != nil {
			return
		}

		//cancelSubOrderSignalBotRequest := new(model.CancelSubOrderSignalBotRequest)
		//cancelSubOrderSignalBotRequest.InstId = "SOL-USDT-SWAP"
		//cancelSubOrderSignalBotRequest.AlgoId = "1811656338578079744"
		//cancelSubOrderSignalBotRequest.SignalOrdId = "1799976084486750208"
		//
		//_, err = app.CancelSubOrderSignalBot(api, cancelSubOrderSignalBotRequest)
		//if err != nil {
		//	return
		//}

		closePositionSignalBotRequest := new(model.ClosePositionSignalBotRequest)
		closePositionSignalBotRequest.InstId = "SOL-USDT-SWAP"
		//closePositionSignalBotRequest.AlgoId = "1811656338578079744"
		//closePositionSignalBotRequest.AlgoId = "1811814804349255680"

		//_, err = app.ClosePositionSignalBot(api, closePositionSignalBotRequest)
		if err != nil {
			return
		}

		getSubOrdersSignalBotRequest := new(model.GetSubOrdersSignalBotRequest)
		getSubOrdersSignalBotRequest.AlgoId = "1817262238944722944"
		getSubOrdersSignalBotRequest.AlgoOrdType = "contract"
		getSubOrdersSignalBotRequest.SignalOrdId = "1799976084486750208"
		//getSubOrdersSignalBotRequest.AlgoId = "1811814804349255680"

		details, err := app.GetSubOrderSignalBot(api, getSubOrdersSignalBotRequest)
		if err != nil {
			return
		}
		logger.Infof("details Data: %v", details)

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
