package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
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

type CalcPriceData struct {
	demoStep      float64
	demoMinAmount float64
	demoPrecision int
	prodStep      float64
	prodMinAmount float64
	prodPrecision int
}

const (
	ADA  = "ADA"
	ATOM = "ATOM"
	AVAX = "AVAX"
	APT  = "APT"

	BTC = "BTC"
	BCH = "BCH"
	BNB = "BNB"

	CRO = "CRO"

	ETH = "ETH"
	ETC = "ETC"

	DOGE = "DOGE"
	DOT  = "DOT"

	FIL = "FIL"

	ICP = "ICP"
	IMX = "IMX"

	LTC  = "LTC"
	NEAR = "NEAR"

	SOL  = "SOL"
	SHIB = "SHIB"

	TON = "TON"

	USD  = "USD"
	USDT = "USDT"
	UNI  = "UNI"

	VET = "VET"

	XRP = "XRP"
	XLM = "XLM"
)

// for calcPx
// https://www.okx.com/ru/trade-market/info/swap
var calcPriceData = map[string]CalcPriceData{
	ETH: {demoStep: 0.1, demoMinAmount: 0.001, demoPrecision: 1, prodStep: 0.1, prodMinAmount: 0.01, prodPrecision: 1},
	XRP: {demoStep: 0.1, demoMinAmount: 10, demoPrecision: 1, prodStep: 0.1, prodMinAmount: 10, prodPrecision: 1},
	SOL: {demoStep: 0.01, demoMinAmount: 0.001, demoPrecision: 2, prodStep: 0.01, prodMinAmount: 0.01, prodPrecision: 2},
	ADA: {demoStep: 0.1, demoMinAmount: 0.1, demoPrecision: 1, prodStep: 0.1, prodMinAmount: 10, prodPrecision: 1},
}

type DealStart struct {
	DealSignal    models.Signal
	DealBot       models.Bot
	DealDirection models.DealDirection
}

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
			err := closeDeal(&deal, &bot, signal.NameToken)
			if err != nil {
				logger.Errorf("Clode deal error: %v", err)
			}
		} else {
			logger.Errorf("There is no deal for the bot=%v and it cannot be closed", bot.ID)
		}
		time.Sleep(time.Second * 1)
	}
}

func openDeal(bot *models.Bot, currencyName string, direction models.DealDirection) error {
	deal := new(models.Deal)
	deal.StartTime = time.Now()
	deal.Direction = direction
	beforeAvailAmount, beforeFrozenAmount, err := getAmount(bot.UserId, bot.OkxBotId, bot.IsProduction)
	if err != nil {
		return err
	}
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

		if px > 0 {
			time.Sleep(time.Second * 3)
			afterAvailAmount, afterFrozenAmount, err := getAmount(bot.UserId, bot.OkxBotId, bot.IsProduction)
			if err != nil {
				return err
			}
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

func closeDeal(deal *models.Deal, bot *models.Bot, currencyName string) error {
	beforeAvailAmount, beforeFrozenAmount, err := getAmount(bot.UserId, bot.OkxBotId, bot.IsProduction)
	if err != nil {
		return err
	}
	if beforeFrozenAmount > 0 {
		result := closeOrder(bot.UserId, currencyName, bot.OkxBotId, bot.IsProduction)
		if result {
			time.Sleep(time.Second * 4)
			afterAvailAmount, afterFrozenAmount, err := getAmount(bot.UserId, bot.OkxBotId, bot.IsProduction)
			if err != nil {
				return err
			}
			var times int = 10
			for afterFrozenAmount > 0 {
				time.Sleep(time.Second * 4)
				afterAvailAmount, afterFrozenAmount, err = getAmount(bot.UserId, bot.OkxBotId, bot.IsProduction)
				times--
				if times == 0 {
					break
				}
			}
			if afterFrozenAmount == 0 {
				diffAmount := afterAvailAmount - beforeAvailAmount
				bot.CurrentAmount = afterAvailAmount
				bot.Status = models.Waiting

				bot.Update()
				deal.FinishDbSave(diffAmount)
			} else {
				return errors.New("Error: afterFrozenAmount not zero")
			}
		} else {
			logger.Errorf("An order on a crypto exchange cannot be closed for a bot=%v and a deal=%v", bot.ID, deal.ID)
		}
	}
	return nil
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
	float64Sz := calcPx(bot.UserId, currencyName, beforeAvailAmount*lever, percent, bot.IsProduction)
	stringSz := strconv.FormatFloat(float64Sz, 'f', 2, 64)
	logger.Infof("Opening order for bot with id: %v and calcPx: %v", bot.ID, stringSz)
	operationCode, err := OkxPlaceSubOrder(bot.UserId, currencyName+"-"+BASE_CURRENCY+"-SWAP", bot.OkxBotId, stringSz, direction, bot.IsProduction)
	if err != nil {
		return "", 0, err
	}
	time.Sleep(2 * time.Second)
	logger.Infof("Code for OkxPlaceSubOrder is %s", operationCode)

	return OkxGetSubOrderSignalBot(bot.UserId, bot.OkxBotId, bot.IsProduction), uint(float64Sz), nil
}

func closeOrder(userId uint, currencyName string, algoId string, isProduction bool) bool {
	err := OkxClosePositionSignalBot(userId, currencyName, algoId, isProduction)
	if err != nil {
		logger.Errorf("Error in OkxClosePositionSignalBot: %v", err)
		return false
	}
	time.Sleep(4 * time.Second)
	return true
}

func getAmount(userId uint, algoId string, isProduction bool) (availBal float64, frozenBal float64, err error) {
	signalBotData := OkxGetSignalBot(userId, algoId, isProduction)
	if signalBotData == nil {
		logger.Errorf("Zero signalBot data")
		return 0, 0, errors.New("Zero signalBot data")
	}
	if signalBotData.AvailBal == "" {
		logger.Errorf("Empty AvailBal")
		return 0, 0, errors.New("Empty AvailBal")
	}
	availBal, err = strconv.ParseFloat(signalBotData.AvailBal, 64)
	frozenBal, err = strconv.ParseFloat(signalBotData.FrozenBal, 64)
	if err != nil {
		logger.Errorf("Error during ParseFloat in GetActiveSignalBot: %v", err)
		return 0, 0, err
	}
	return availBal, frozenBal, nil
}

func calcPx(userId uint, symbol string, amount float64, percent float64, isProduction bool) float64 {
	ticker := OkxGetTicker(userId, symbol, isProduction)
	price := ticker.Last
	calcData := calcPriceData[symbol]
	if isProduction {
		return Round(calcData.prodStep*percent*amount/(calcData.prodMinAmount*price*100), calcData.prodPrecision)
	}
	return Round(calcData.demoStep*percent*amount/(calcData.demoMinAmount*price*100), calcData.demoPrecision)
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
	u.Respond(w, u.Message(true, "The checking was finished"))
}
