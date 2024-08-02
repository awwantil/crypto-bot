package app

import (
	"github.com/sirupsen/logrus"
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx"
	"okx-bot/exchange/okx/futures"
	"okx-bot/exchange/options"
	"okx-bot/frontend-service/models"
	"strconv"
	"strings"
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

func CreateOrder(api *futures.PrvApi, pairName string, posSide string) (id string, err error) {
	//https://www.okx.com/docs-v5/en/#order-book-trading-trade-post-place-order
	orderRequest := new(model.PlaceOrderRequest)
	orderRequest.InstId = pairName + "-SWAP"
	orderRequest.TdMode = "isolated"
	orderRequest.Side = "buy"
	orderRequest.OrdType = "market"
	orderRequest.Sz = posSide
	//orderRequest.PxUsd = "12"

	newOrder, data, err := api.Isolated.PlaceOrder(*orderRequest)
	if err != nil {
		logger.Errorf("Error place order: %v, data: %v", err, string(data))
		return "", err
	}
	orderId := newOrder.Id
	logger.Info("ordId = ", orderId)

	return orderId, nil
}

func EndDeal(api *futures.PrvApi, pairName string, orderId string, posSide string) error {
	closePositionsRequest := new(model.ClosePositionsRequest)
	closePositionsRequest.InstId = pairName + "-SWAP"
	closePositionsRequest.ClOrdId = orderId
	respModel, data, err := api.Isolated.ClosePositions(closePositionsRequest)
	if err != nil {
		logger.Errorf("Error close order: %v data: %v", err, data)
		return err
	}
	logger.Info("CLose order respModel: ", respModel)
	return nil
}

func GetOkxAmount(api *futures.PrvApi, userId uint, currencyName string) (float64, float64) {
	requestBalance := new(model.BalanceRequest)
	requestBalance.CCY = currencyName
	resp, data, err := api.GetAccountBalance(*requestBalance)
	if err != nil {
		logger.Errorf("Error in GetOkxAmount: %v", err)
	}
	logger.Infoln("data", string(data))
	if len(resp.Details) > 0 {
		availBalance, _ := strconv.ParseFloat(strings.TrimSpace(resp.Details[0].AvailBal), 64)
		frozenBalance, _ := strconv.ParseFloat(strings.TrimSpace(resp.Details[0].FrozenBal), 64)
		return availBalance, frozenBalance
	}

	return 0.0, 0.0
}
