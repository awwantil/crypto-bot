package controllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
	"sync"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.signal-controllers",
	})
	longNamesArray  = []string{"Long", "long"}
	shortNamesArray = []string{"Short", "short"}
)

const (
	SELL        string = "sell"
	BUY         string = "buy"
	ZERO_AMOUNT string = "0"
	JOBS_NUMBER int    = 10
)

var ReceiveSignal = func(w http.ResponseWriter, r *http.Request) {

	signal := &models.TradingViewSignalReceive{}

	err := json.NewDecoder(r.Body).Decode(signal)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	signal.Save()

	if isLongPosition(signal) {
		if signal.Action == BUY {
			err := startDeal(signal.SignalToken, models.Long)
			if err != nil {
				u.Respond(w, u.Message(false, err.Error()))
				return
			}
		}
		if signal.Action == SELL {
			endDeal(signal.SignalToken)
		}
	}

	if isShortPosition(signal) {
		if signal.Action == SELL {
			err := startDeal(signal.SignalToken, models.Short)
			if err != nil {
				u.Respond(w, u.Message(false, err.Error()))
				return
			}
		}
		if signal.Action == BUY {
			endDeal(signal.SignalToken)
		}
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

	resp := signal.Create()
	u.Respond(w, resp)
}

var GetAllSignals = func(w http.ResponseWriter, r *http.Request) {

	signals := models.GetAllSignals()
	resp := u.Message(true, "success")
	resp["signals"] = signals

	logger.Infoln("resp", resp)
	u.Respond(w, resp)
}

var CheckOkx = func(w http.ResponseWriter, r *http.Request) {
	u.Respond(w, u.Message(true, "The checking was finished"))
}

func startDeal(signalCode string, direction models.DealDirection) error {
	signal, err := models.FindSignalByCode(signalCode)
	if err != nil {
		logger.Errorf("Error in start deal: %v", err)
		return err
	}

	bots := models.GetBots(signalCode)
	if len(bots) > 0 {
		var wg sync.WaitGroup
		semaphore := NewSemaphore(JOBS_NUMBER)
		for _, bot := range bots {
			deal := models.FindByStatus(bot.ID, models.DealStarted)
			if deal.ID == 0 {
				logger.Infof("starting bot's deal with bot id %d and direction: %v", bot.ID, direction)
				wg.Add(1)
				dealStart := DealStart{
					*signal, bot, direction,
				}

				go func(dealStarting *DealStart) {
					semaphore.Acquire()
					defer wg.Done()
					defer semaphore.Release()

					err := dealStarting.openDeal()
					if err != nil {
						logger.Errorf("start deal error: %v", err)
						dealStarting.saveError(err)
					}
				}(&dealStart)
			} else {
				logger.Errorf("There is already a deal=%v for the bot=%v", deal.ID, bot.ID)
			}
		}
		wg.Wait()
	}
	return nil
}

func endDeal(signalCode string) {
	signal, err := models.FindSignalByCode(signalCode)
	if err != nil {
		logger.Errorf("Error in endDeal: %v", err)
		return
	}

	bots := models.GetBots(signalCode)
	if len(bots) > 0 {
		var wg sync.WaitGroup
		semaphore := NewSemaphore(JOBS_NUMBER)
		for _, bot := range bots {
			deal := models.FindByStatus(bot.ID, models.DealStarted)
			if deal.ID > 0 {
				wg.Add(1)
				dealFinish := DealFinish{
					*signal, bot, deal,
				}
				go func(dealFinishing *DealFinish) {
					semaphore.Acquire()
					defer wg.Done()
					defer semaphore.Release()

					err := dealFinishing.closeDeal()
					if err != nil {
						logger.Errorf("Clode deal error: %v", err)
						dealFinishing.saveError(err)
					}
				}(&dealFinish)
			} else {
				logger.Errorf("There is no deal for the bot=%v and it cannot be closed", bot.ID)
			}
		}
		wg.Wait()
	}
}

func isLongPosition(signal *models.TradingViewSignalReceive) bool {
	signalId := signal.Id
	for _, currentName := range longNamesArray {
		if findSubstring(signalId, currentName) {
			return true
		}
	}
	if signal.Action == BUY && signal.MarketPositionSize != ZERO_AMOUNT && signal.PrevMarketPositionSize == ZERO_AMOUNT {
		return true
	}
	if signal.Action == SELL && signal.MarketPositionSize == ZERO_AMOUNT && signal.PrevMarketPositionSize != ZERO_AMOUNT {
		return true
	}
	return false
}

func isShortPosition(signal *models.TradingViewSignalReceive) bool {
	signalId := signal.Id
	for _, currentName := range shortNamesArray {
		if findSubstring(signalId, currentName) {
			return true
		}
	}
	if signal.Action == SELL && signal.MarketPositionSize != ZERO_AMOUNT && signal.PrevMarketPositionSize == ZERO_AMOUNT {
		return true
	}
	if signal.Action == BUY && signal.MarketPositionSize == ZERO_AMOUNT && signal.PrevMarketPositionSize != ZERO_AMOUNT {
		return true
	}
	return false
}

func findSubstring(str string, match string) bool {
	if len(str) < len(match) {
		return false
	}
	for i := 0; i <= len(str)-len(match); i++ {
		subStr := str[i : i+len(match)]
		if subStr == match {
			return true
		}
	}
	return false
}
