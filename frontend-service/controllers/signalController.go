package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
	"okx-bot/frontend-service/app"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
	"strconv"
	"time"
)

//https://www.okx.com/ru/trade-market/info/swap

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.signal-controllers",
	})
	longNamesArray  = []string{"Long", "long"}
	shortNamesArray = []string{"Short", "short"}
)

const (
	SELL            string  = "sell"
	BUY             string  = "buy"
	LONG            string  = "long"
	SHORT           string  = "short"
	FLAT            string  = "flat"
	ZERO_AMOUNT     string  = "0"
	BASE_CURRENCY           = "USDT"
	DEFAULT_PERCENT float64 = 60
	DEFAULT_LEVER   float64 = 3
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

func startDeal(signalCode string, direction models.DealDirection) error {
	signal, err := models.FindSignalByCode(signalCode)
	if err != nil {
		logger.Errorf("Error in start deal: %v", err)
		return err
	}
	bots := models.GetBots(signalCode)
	for _, bot := range bots {
		deal := models.FindByStatus(bot.ID, models.DealStarted)
		if deal.ID == 0 {
			logger.Infof("starting bot's deal with bot id %d and direction: %v", bot.ID, direction)
			err := openDeal(&bot, signal.NameToken, direction)
			logger.Errorf("start deal error: %v", err)
		} else {
			logger.Errorf("There is already a deal=%v for the bot=%v", deal.ID, bot.ID)
		}
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
	for _, bot := range bots {
		deal := models.FindByStatus(bot.ID, models.DealStarted)
		if deal.ID > 0 {
			closeDeal(&deal, &bot, signal.NameToken)
		} else {
			logger.Errorf("There is no deal for the bot=%v and it cannot be closed", bot.ID)
		}
	}
}

func openDeal(bot *models.Bot, currencyName string, direction models.DealDirection) error {
	deal := new(models.Deal)
	deal.StartTime = time.Now()
	deal.Direction = direction
	beforeAvailAmount, beforeFrozenAmount := getAmount(bot.UserId, bot.OkxBotId)
	logger.Infof("OpenDeal before available amount: %v, frozen amount: %v", beforeAvailAmount, beforeFrozenAmount)
	if beforeFrozenAmount > 0 {
		return errors.New("The deal is already exist")
	}

	if beforeAvailAmount > 0 {
		order, px, err := openOrder(bot, currencyName, beforeAvailAmount, direction)
		if err != nil {
			return err
		}
		deal.OrderId = order

		if deal.OrderId != "" {
			afterAvailAmount, afterFrozenAmount := getAmount(bot.UserId, bot.OkxBotId)
			logger.Infof("After available amount: %v, frozen amount: %v", afterAvailAmount, afterFrozenAmount)
			diffAmount := beforeAvailAmount - afterAvailAmount
			deal.StartAmount = diffAmount
			deal.Status = models.DealStarted
			deal.StartDbSave(bot.ID, diffAmount)
			bot.CurrentAmount = afterAvailAmount
			bot.Status = models.MakingDeal
			bot.PosSide = px
			bot.Update()
		} else {
			strErr := fmt.Sprintf("An order cannot be created for a bot=%v on a crypto exchange", bot.ID)
			return errors.New(strErr)
		}
	} else {
		strErr := fmt.Sprintf("For openDeal beforeAvailAmount equals zero")
		logger.Errorf("For openDeal beforeAvailAmount equals zero")
		return errors.New(strErr)
	}
	return nil
}

func closeDeal(deal *models.Deal, bot *models.Bot, currencyName string) {
	beforeAvailAmount, beforeFrozenAmount := getAmount(bot.UserId, bot.OkxBotId)
	if beforeFrozenAmount > 0 {
		result := closeOrder(bot.UserId, currencyName, bot.OkxBotId)
		if result {
			afterAvailAmount, afterFrozenAmount := getAmount(bot.UserId, bot.OkxBotId)
			if afterFrozenAmount == 0 {
				diffAmount := afterAvailAmount - beforeAvailAmount
				bot.CurrentAmount = afterAvailAmount
				bot.Status = models.Waiting

				bot.Update()
				deal.FinishDbSave(diffAmount)
			}
		} else {
			logger.Errorf("An order on a crypto exchange cannot be closed for a bot=%v and a deal=%v", bot.ID, deal.ID)
		}
	}
}

func openOrder(bot *models.Bot, currencyName string, beforeAvailAmount float64, direction models.DealDirection) (string, uint, error) {
	percent := bot.DealsPercent
	lever := bot.Lever

	if percent == 0 {
		percent = DEFAULT_PERCENT
	}
	if lever == 0 {
		lever = DEFAULT_LEVER
	}
	logger.Infof("amount: %v", beforeAvailAmount*lever)
	float64Sz := calcPx(bot.UserId, currencyName, beforeAvailAmount*lever, percent)
	stringSz := strconv.FormatFloat(float64Sz, 'f', 2, 64)
	logger.Infof("calcPx: %v", stringSz)
	operationCode, err := OkxPlaceSubOrder(bot.UserId, currencyName+"-"+BASE_CURRENCY+"-SWAP", bot.OkxBotId, stringSz, direction)
	if err != nil {
		return "", 0, err
	}
	time.Sleep(2 * time.Second)
	logger.Infof("Code for OkxPlaceSubOrder is %s", operationCode)

	return OkxGetSubOrderSignalBot(bot.UserId, bot.OkxBotId), uint(float64Sz), nil
}

func closeOrder(userId uint, currencyName string, algoId string) bool {
	err := OkxClosePositionSignalBot(userId, currencyName, algoId)
	if err != nil {
		logger.Errorf("Error in OkxClosePositionSignalBot: %v", err)
		return false
	}
	return true
}

func getAmount(userId uint, algoId string) (float64, float64) {
	signalBotData := OkxGetSignalBot(userId, algoId)
	if signalBotData == nil {
		logger.Errorf("Error in request signal bot")
		return 0, 0
	}
	if signalBotData.AvailBal == "" {
		logger.Errorf("Error in request to amounts")
		return 0, 0
	}
	availBal, err := strconv.ParseFloat(signalBotData.AvailBal, 64)
	frozenBal, err := strconv.ParseFloat(signalBotData.FrozenBal, 64)
	if err != nil {
		logger.Errorf("Error in GetActiveSignalBot: %v", err)
		return 0, 0
	}
	return availBal, frozenBal
}

func getBaseAmount(userId uint, currencyName string) (float64, float64) {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return 0, 0
	}
	return app.GetOkxAmount(api, userId, currencyName)
}

func calcPx(userId uint, symbol string, amount float64, percent float64) float64 {
	ticker := OkxGetTicker(userId, symbol)
	price := ticker.Last
	if symbol == "ETH" {
		return Round(percent*amount/(price), 1)
	}
	if symbol == "XRP" {
		return Round(percent*amount/(price*10000), 1)
	}
	return Round(percent*amount/(price*10), 2)
}

func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
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

var CheckOkx = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint)

	activeSignalBot := OkxGetActiveSignalBot(user, "1843697369594986496")
	if activeSignalBot != nil {
		logger.Info("Active signalBot: ", activeSignalBot)
		logger.Info("Active AvailBal: ", activeSignalBot.AvailBal)
		logger.Info("Active FrozenBal: ", activeSignalBot.FrozenBal)
	}

	signalBot := OkxGetSignalBot(user, "1843697369594986496")
	if signalBot != nil {
		logger.Info("signalBot: ", signalBot)
		logger.Info("AvailBal: ", signalBot.AvailBal)
		logger.Info("FrozenBal: ", signalBot.FrozenBal)
	}

	u.Respond(w, u.Message(true, "The checking was finished"))
}
