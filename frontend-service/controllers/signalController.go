package controllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
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
			deal := new(models.Deal)
			deal.StartTime = time.Now()
			startAmount := float64(23)
			deal.Start(bot.ID, startAmount)
		} else {
			endAmount := float64(12)
			deal.Failure(endAmount)
		}
	}
}

func endDeal(signalCode string) {
	bots := models.GetBots(signalCode)
	for _, bot := range bots {
		logger.Infof("end bot's deal with id %d", bot.ID)
		deal := models.FindByStatus(bot.ID, models.DealStarted)
		if deal.ID > 0 {
			endAmount := float64(23)
			deal.Finish(endAmount)
		}
	}
}
