package controllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"okx-bot/frontend-service/app"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
	"time"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.signal-controllers",
	})
)

const (
	SELL string = "sell"
	BUY  string = "buy"
)

var ReceiveSignal = func(w http.ResponseWriter, r *http.Request) {

	signal := &models.TradingViewSignalReceive{}

	err := json.NewDecoder(r.Body).Decode(signal)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	signal.Save()

	if signal.Action == BUY {
		startDeal(signal.SignalToken)
	}
	if signal.Action == SELL {
		endDeal(signal.SignalToken)
	}

	u.Respond(w, u.Message(true, "The signal was received"))
}

var CreateSignal = func(w http.ResponseWriter, r *http.Request) {

	signal := &models.Signal{}

	err := json.NewDecoder(r.Body).Decode(signal)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	resp := signal.Create(signal.NameToken, signal.TimeInterval)
	u.Respond(w, resp)
}

func startDeal(signalCode string) {
	bots := models.GetBots(signalCode)
	for _, bot := range bots {
		logger.Infof("start bot's deal with id %d", bot.ID)
		deal := models.FindByStatus(bot.ID, models.DealStarted)
		if deal.ID == 0 {
			//start deal
			deal := new(models.Deal)
			deal.StartTime = time.Now()
			beforeAvailAmount, beforeFrozenAmount := getAmount(bot.UserId, "USDT")
			logger.Infof("Before available amount: %d, frozen amount: %d", beforeAvailAmount, beforeFrozenAmount)

			deal.OrderId = openOrder(bot.UserId, "SOL")

			if deal.OrderId != "" {
				afterAvailAmount, afterFrozenAmount := getAmount(bot.UserId, "USDT")
				logger.Infof("After available amount: %d, frozen amount: %d", afterAvailAmount, afterFrozenAmount)
				diffAmount := beforeAvailAmount - afterAvailAmount
				deal.StartAmount = diffAmount
				deal.StartDbSave(bot.ID, diffAmount)
				bot.CurrentAmount = diffAmount
				bot.Update()
			}
		} else {
			availAmount, frozenAmount := getAmount(bot.UserId, "USDT")
			logger.Infof("Available amount: %d, frozen amount: %d", availAmount, frozenAmount)
			//deal.Failure(endAmount)
			//bot.Status = models.Waiting
			//bot.Update("status")
		}
	}
}

func endDeal(signalCode string) {
	bots := models.GetBots(signalCode)
	for _, bot := range bots {
		logger.Infof("end bot's deal with id %d", bot.ID)
		deal := models.FindByStatus(bot.ID, models.DealStarted)
		if deal.ID > 0 {
			beforeAvailAmount, beforeFrozenAmount := getAmount(bot.UserId, "USDT")
			logger.Infof("Before available amount: %d, frozen amount: %d", beforeAvailAmount, beforeFrozenAmount)

			result := closeOrder(bot.UserId, "SOL", deal.OrderId)
			if result {
				afterAvailAmount, afterFrozenAmount := getAmount(bot.UserId, "USDT")
				logger.Infof("After available amount: %d, frozen amount: %d", afterAvailAmount, afterFrozenAmount)
				diffAmount := afterAvailAmount - beforeAvailAmount
				bot.CurrentAmount = diffAmount
				bot.Update()
				deal.FinishDbSave(diffAmount)
			}
		}
	}
}

func openOrder(userId uint, currencyName string) (orderId string) {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return ""
	}

	orderId, err = app.CreateOrder(api, currencyName+"-USDT", "2")
	if err != nil {
		return ""
	}

	return orderId
}

func closeOrder(userId uint, currencyName string, orderId string) bool {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return false
	}

	err = app.EndDeal(api, currencyName+"-USDT", orderId, "2")
	if err != nil {
		return false
	}

	return true
}

func getAmount(userId uint, currencyName string) (float64, float64) {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return 0, 0
	}
	return app.GetOkxAmount(api, userId, currencyName)
}

var CheckOkx = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint)

	availAmount, frozenAmount := getAmount(user, "USDT")
	logger.Infof("USDT available amount: %v, frozen amount: %v", availAmount, frozenAmount)
	availAmount, frozenAmount = getAmount(user, "SOL")
	logger.Infof("SOL available amount: %v, frozen amount: %v", availAmount, frozenAmount)

	u.Respond(w, u.Message(true, "The checking was finished"))
}
