package controllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"okx-bot/exchange/model"
	"okx-bot/frontend-service/app"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
	"strconv"
	"strings"
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
			availAmount, frozenAmount := getOkxAmount(bot.UserId, "USDT")
			logger.Infof("Available amount: %d, frozen amount: %d", availAmount, frozenAmount)
			deal.StartDbSave(bot.ID, bot.CurrentAmount)
			//bot.CurrentAmount = startAmount
			//bot.Update("current_amount")
		} else {
			availAmount, frozenAmount := getOkxAmount(bot.UserId, "USDT")
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
		//stop deal
		availAmount, frozenAmount := getOkxAmount(bot.UserId, "USDT")
		logger.Infof("Available amount: %d, frozen amount: %d", availAmount, frozenAmount)
		endedAmount := availAmount
		bot.CurrentAmount = endedAmount
		bot.Update()
		if deal.ID > 0 {
			endAmount := endedAmount
			deal.FinishDbSave(endAmount)
		}
	}
}

func getOkxAmount(userId uint, currencyName string) (float64, float64) {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in getOkxAmount: %v", err)
		return 0, 0
	}
	logger.Info("api", api.GetName())
	requestBalance := new(model.BalanceRequest)
	requestBalance.CCY = currencyName
	resp, data, err := api.GetAccountBalance(*requestBalance)
	if err != nil {
		logger.Errorf("Error in getOkxAmount: %v", err)
	}
	logger.Infoln("data", data)
	logger.Infoln("resp", resp)
	if len(resp.Details) > 0 {
		availBalance, _ := strconv.ParseFloat(strings.TrimSpace(resp.Details[0].AvailBal), 64)
		frozenBalance, _ := strconv.ParseFloat(strings.TrimSpace(resp.Details[0].FrozenBal), 64)
		return availBalance, frozenBalance
	}

	return 0.0, 0.0
}

var CheckOkx = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint)

	availAmount, frozenAmount := getOkxAmount(user, "USDT")
	logger.Infof("USDT available amount: %v, frozen amount: %v", availAmount, frozenAmount)
	availAmount, frozenAmount = getOkxAmount(user, "SOL")
	logger.Infof("SOL available amount: %v, frozen amount: %v", availAmount, frozenAmount)

	u.Respond(w, u.Message(true, "The checking was finished"))
}
