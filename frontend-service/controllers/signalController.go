package controllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"okx-bot/frontend-service/app"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
	"strconv"
	"time"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.signal-controllers",
	})
)

const (
	SELL          string = "sell"
	BUY           string = "buy"
	BASE_CURRENCY        = "USDT"
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

	resp := signal.Create(signal.StrategyName, signal.NameToken, signal.TimeInterval)
	u.Respond(w, resp)
}

var GetAllSignals = func(w http.ResponseWriter, r *http.Request) {

	signals := models.GetAllSignals()
	resp := u.Message(true, "success")
	resp["signals"] = signals

	logger.Infoln("resp", resp)
	u.Respond(w, resp)
}

func startDeal(signalCode string) {
	signal, err := models.FindSignalByCode(signalCode)
	if err != nil {
		logger.Errorf("Error in startDeal: %v", err)
		return
	}
	bots := models.GetBots(signalCode)
	for _, bot := range bots {
		logger.Infof("start bot's deal with id %d", bot.ID)
		deal := models.FindByStatus(bot.ID, models.DealStarted)
		if deal.ID == 0 {
			openDeal(&bot, signal.NameToken)
		} else {
			logger.Errorf("There is already a deal=%v for the bot=%v", deal.ID, bot.ID)
		}
	}
}

func endDeal(signalCode string) {
	signal, err := models.FindSignalByCode(signalCode)
	if err != nil {
		logger.Errorf("Error in startDeal: %v", err)
		return
	}
	bots := models.GetBots(signalCode)
	for _, bot := range bots {
		logger.Infof("end bot's deal with id %d", bot.ID)
		deal := models.FindByStatus(bot.ID, models.DealStarted)
		if deal.ID > 0 {
			closeDeal(&deal, &bot, signal.NameToken)
		} else {
			logger.Errorf("There is no deal for the bot=%v and it cannot be closed", bot.ID)
		}
	}
}

func openDeal(bot *models.Bot, currencyName string) {
	deal := new(models.Deal)
	deal.StartTime = time.Now()
	beforeAvailAmount, beforeFrozenAmount := getAmount(bot.UserId, BASE_CURRENCY)
	logger.Infof("Before available amount: %d, frozen amount: %d", beforeAvailAmount, beforeFrozenAmount)

	posSide := "2"
	if bot.PosSide > 0 {
		posSide = strconv.FormatUint(uint64(bot.PosSide), 10)
	}
	deal.OrderId = openOrder(bot.UserId, currencyName, posSide)

	if deal.OrderId != "" {
		afterAvailAmount, afterFrozenAmount := getAmount(bot.UserId, BASE_CURRENCY)
		logger.Infof("After available amount: %d, frozen amount: %d", afterAvailAmount, afterFrozenAmount)
		diffAmount := beforeAvailAmount - afterAvailAmount
		deal.StartAmount = diffAmount
		deal.StartDbSave(bot.ID, diffAmount)
		bot.CurrentAmount = diffAmount
		bot.Update()
	} else {
		logger.Errorf("An order cannot be created for a bot=%v on a crypto exchange", bot.ID)
	}
}

func closeDeal(deal *models.Deal, bot *models.Bot, currencyName string) {
	beforeAvailAmount, beforeFrozenAmount := getAmount(bot.UserId, BASE_CURRENCY)
	logger.Infof("Before available amount: %d, frozen amount: %d", beforeAvailAmount, beforeFrozenAmount)
	result := closeOrder(bot.UserId, currencyName, deal.OrderId)
	if result {
		afterAvailAmount, afterFrozenAmount := getAmount(bot.UserId, BASE_CURRENCY)
		logger.Infof("After available amount: %d, frozen amount: %d", afterAvailAmount, afterFrozenAmount)
		diffAmount := afterAvailAmount - beforeAvailAmount
		bot.CurrentAmount = diffAmount
		bot.Update()
		deal.FinishDbSave(diffAmount)
	} else {
		logger.Errorf("An order on a crypto exchange cannot be closed for a bot=%v and a deal=%v", bot.ID, deal.ID)
	}
}

func openOrder(userId uint, currencyName string, posSide string) (orderId string) {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return ""
	}

	orderId, err = app.CreateOrder(api, currencyName+"-"+BASE_CURRENCY, posSide)
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

	err = app.EndDeal(api, currencyName+"-"+BASE_CURRENCY, orderId, "2")
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

	availAmount, frozenAmount := getAmount(user, BASE_CURRENCY)
	logger.Infof("USDT available amount: %v, frozen amount: %v", availAmount, frozenAmount)
	availAmount, frozenAmount = getAmount(user, "SOL")
	logger.Infof("SOL available amount: %v, frozen amount: %v", availAmount, frozenAmount)

	u.Respond(w, u.Message(true, "The checking was finished"))
}
