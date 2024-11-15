package controllers

import (
	"errors"
	"fmt"
	"math"
	"okx-bot/frontend-service/models"
	"strconv"
	"time"
)

const (
	BASE_CURRENCY           = "USDT"
	DEFAULT_PERCENT float64 = 60
	DEFAULT_LEVER   float64 = 3
)

type DealStart struct {
	DealSignal    models.Signal
	DealBot       models.Bot
	DealDirection models.DealDirection
}

type DealFinish struct {
	DealSignal models.Signal
	DealBot    models.Bot
	Deal       models.Deal
}

func (dealStart *DealStart) openDeal() error {
	direction := dealStart.DealDirection
	bot := dealStart.DealBot

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
		order, px, err := dealStart.openOrder(beforeAvailAmount)
		if err != nil {
			return err
		}
		deal.OrderId = order

		if px > 0 || order != "" {
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
			bot.PosSide = uint(px)
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

func (dealFinish *DealFinish) closeDeal() error {
	bot := dealFinish.DealBot
	deal := dealFinish.Deal
	beforeAvailAmount, beforeFrozenAmount, err := getAmount(bot.UserId, bot.OkxBotId, bot.IsProduction)
	if err != nil {
		return err
	}
	if beforeFrozenAmount > 0 {
		result := dealFinish.closeOrder()
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
			strErr := fmt.Sprintf("An order on a crypto exchange cannot be closed for a bot=%v and a deal=%v", bot.ID, deal.ID)
			return errors.New(strErr)
		}
	}
	return nil
}

func getAmount(userId uint, algoId string, isProduction bool) (availBal float64, frozenBal float64, err error) {
	signalBotData := OkxGetSignalBot(userId, algoId, isProduction)
	if signalBotData == nil {
		strErr := fmt.Sprintf("Zero signalBot data")
		return 0, 0, errors.New(strErr)
	}
	if signalBotData.AvailBal == "" {
		strErr := fmt.Sprintf("Empty AvailBal")
		return 0, 0, errors.New(strErr)
	}
	availBal, err = strconv.ParseFloat(signalBotData.AvailBal, 64)
	frozenBal, err = strconv.ParseFloat(signalBotData.FrozenBal, 64)
	if err != nil {
		return 0, 0, err
	}
	return availBal, frozenBal, nil
}

func (dealStart *DealStart) openOrder(beforeAvailAmount float64) (string, float64, error) {
	bot := dealStart.DealBot
	currencyName := dealStart.DealSignal.NameToken
	direction := dealStart.DealDirection

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

	return OkxGetSubOrderSignalBot(bot.UserId, bot.OkxBotId, bot.IsProduction), float64Sz, nil
}

func (dealFinish *DealFinish) closeOrder() bool {
	userId := dealFinish.DealBot.UserId
	currencyName := dealFinish.DealSignal.NameToken
	algoId := dealFinish.DealBot.OkxBotId
	isProduction := dealFinish.DealBot.IsProduction

	err := OkxClosePositionSignalBot(userId, currencyName, algoId, isProduction)
	if err != nil {
		logger.Errorf("Error in OkxClosePositionSignalBot: %v", err)
		return false
	}
	time.Sleep(4 * time.Second)
	return true
}

func calcPx(userId uint, symbol string, amount float64, percent float64, isProduction bool) float64 {
	ticker := OkxGetTicker(userId, symbol, isProduction)
	price := ticker.Last
	calcData := models.PriceData[symbol]
	if isProduction {
		return Round(calcData.ProdStep*percent*amount/(calcData.ProdMinAmount*price*100), calcData.ProdPrecision)
	}
	return Round(calcData.DemoStep*percent*amount/(calcData.DemoMinAmount*price*100), calcData.DemoPrecision)
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

func (dealStart *DealStart) saveError(err error) {
	botError := models.BotError{
		StartTime:   time.Now(),
		SignalRefer: dealStart.DealSignal.ID,
		OkxBotId:    dealStart.DealBot.OkxBotId,
		Message:     err.Error(),
		IsOpenDeal:  true,
	}
	botError.SaveBotError()
}

func (dealFinish *DealFinish) saveError(err error) {
	botError := models.BotError{
		StartTime:   time.Now(),
		SignalRefer: dealFinish.DealSignal.ID,
		OkxBotId:    dealFinish.DealBot.OkxBotId,
		Message:     err.Error(),
		IsOpenDeal:  false,
	}
	botError.SaveBotError()
}
