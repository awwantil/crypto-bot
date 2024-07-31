package app

import (
	"github.com/sirupsen/logrus"
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx"
	"okx-bot/exchange/okx/futures"
	"okx-bot/exchange/options"
	"okx-bot/frontend-service/models"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.deals",
	})
)

func GetOkxApi(userId uint) (*futures.PrvApi, error) {
	userOkxApi, err := models.GetUserApi(userId)
	if err != nil {
		return nil, err
	}
	OKx := okx.New()
	okx.DefaultHttpCli.SetTimeout(5)

	api := OKx.Futures.NewPrvApi(
		options.WithApiKey(userOkxApi.MainKey),
		options.WithApiSecretKey(userOkxApi.SpecialKey),
		options.WithPassphrase(userOkxApi.Phrase))

	return api, nil
}

func CreateOrder(api *futures.PrvApi, pairName string, posSide string) (id string) {
	//https://www.okx.com/docs-v5/en/#order-book-trading-trade-post-place-order
	orderRequest := new(model.PlaceOrderRequest)
	orderRequest.InstId = pairName + "-SWAP"
	orderRequest.TdMode = "isolated"
	orderRequest.Side = "sell"
	orderRequest.OrdType = "market"
	orderRequest.Sz = "10"
	//orderRequest.PxUsd = "12"

	newOrder, _, err := api.Isolated.PlaceOrder(*orderRequest)
	if err != nil {
		logger.Error("Error place order", err)
	}
	orderId := newOrder.Id
	logger.Info("ordId = ", orderId)

	return orderId
}

func EndDeal(api *futures.PrvApi, pair string, orderId string, posSide string) {
	// ordId = 1341342760297025536
	closePositionsRequest := new(model.ClosePositionsRequest)
	closePositionsRequest.InstId = pair + "-SWAP"
	closePositionsRequest.ClOrdId = orderId
	closePositionsRequest.PosSide = posSide
	respBody, respModel, err := api.Isolated.ClosePositions(closePositionsRequest)
	if err != nil {
		logger.Error("Error close order: ", err)
		return
	}
	logger.Info("CLose order resp: ", respBody)
	logger.Info("CLose order respModel: ", respModel)
}

func getAmount(api *futures.PrvApi) float64 {
	//api.Isolated.
	return 0
}
