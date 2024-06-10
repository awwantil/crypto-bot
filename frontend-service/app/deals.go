package app

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx"
	"okx-bot/exchange/okx/futures"
	"okx-bot/exchange/options"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.deals",
	})
)

var okxApi *futures.PrvApi

func InitOkxApi() {
	logger.Infoln("Starting create Exchange API")

	envParams := make(map[string]string)
	envParams, err := godotenv.Read()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	OKx := okx.New()
	okx.DefaultHttpCli.SetTimeout(5)

	api := OKx.Futures.NewPrvApi(
		options.WithApiKey(envParams["okx_api_key"]),
		options.WithApiSecretKey(envParams["okx_api_secret_key"]),
		options.WithPassphrase(envParams["okx_api_passphrase"]))

	if api == nil {
		logger.Fatal("Error creating Exchange API")
		return
	}

	okxApi = api
	logger.Infoln("Exchange API created")
}

func GetOkxApi() *futures.PrvApi {
	envParams := make(map[string]string)
	envParams, err := godotenv.Read()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	OKx := okx.New()
	okx.DefaultHttpCli.SetTimeout(5)

	api := OKx.Futures.NewPrvApi(
		options.WithApiKey(envParams["okx_api_key"]),
		options.WithApiSecretKey(envParams["okx_api_secret_key"]),
		options.WithPassphrase(envParams["okx_api_passphrase"]))

	if api == nil {
		logger.Fatal("Error creating Exchange OKx API")
		return nil
	}

	return api
}

func StartDeal(pair string, posSide string) (id string) {
	//https://www.okx.com/docs-v5/en/#order-book-trading-trade-post-place-order
	orderRequest := new(model.PlaceOrderRequest)
	orderRequest.InstId = pair + "-SWAP"
	orderRequest.TdMode = "isolated"
	orderRequest.Side = "sell"
	orderRequest.OrdType = "market"
	orderRequest.Sz = "10"
	//orderRequest.PxUsd = "12"
	orderRequest.ClOrdId = "12345"

	newOrder, _, err := okxApi.Isolated.PlaceOrder(*orderRequest)
	if err != nil {
		logger.Error("Error place order", err)
	}
	orderId := newOrder.Id
	logger.Info("ordId = ", orderId)

	return orderId
}

func EndDeal(pair string, ordId string, posSide string) {
	// ordId = 1341342760297025536
	closePositionsRequest := new(model.ClosePositionsRequest)
	closePositionsRequest.InstId = pair + "-SWAP"
	closePositionsRequest.ClOrdId = ordId
	closePositionsRequest.PosSide = posSide
	respBody, respModel, err := okxApi.Isolated.ClosePositions(closePositionsRequest)
	if err != nil {
		logger.Error("Error close order: ", err)
		return
	}
	logger.Info("CLose order resp: ", respBody)
	logger.Info("CLose order respModel: ", respModel)
}

func GetApi() *futures.PrvApi {
	if okxApi == nil {
		logger.Error("Not OKX API connection")
		return nil
	}
	return okxApi
}
