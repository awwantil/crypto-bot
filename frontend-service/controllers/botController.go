package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aidarkhanov/nanoid"
	"net/http"
	"okx-bot/exchange/model"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
	"strconv"
	"time"
)

var CreateBot = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	botRequest := &models.BotCreateRequest{}

	err := json.NewDecoder(r.Body).Decode(botRequest)
	if err != nil {
		logger.Errorf("Error while decoding request body: %s", err.Error())
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	if botRequest.InitialAmount < 50.0 {
		logger.Errorf("The order amount must be greater than 50")
		u.Respond(w, u.Message(false, "The order amount must be greater than 50"))
		return
	}
	if botRequest.Lever < 0 || botRequest.Lever > 50.0 {
		logger.Errorf("The lever is wrong")
		u.Respond(w, u.Message(false, "The lever is wrong"))
		return
	}
	if botRequest.DealsPercent < 0 || botRequest.DealsPercent > 100.0 {
		logger.Errorf("The dealsPercent is wrong")
		u.Respond(w, u.Message(false, "The dealsPercent is wrong"))
		return
	}

	bot := new(models.Bot)
	bot.UserId = user
	bot.Lever = botRequest.Lever
	bot.DealsPercent = botRequest.DealsPercent
	bot.IsProduction = botRequest.IsProduction

	var signal = models.Signal{Code: botRequest.CodeSignalId}
	models.GetDB().Where("code = ?", botRequest.CodeSignalId).First(&signal)

	if signal.NameToken == "" {
		logger.Errorf("The signal is wrong")
		u.Respond(w, u.Message(false, "The signal is wrong"))
		return
	}

	logger.Infoln("signal", signal)

	okxSignalId := getOkxSignalId(user, signal.StrategyName, signal.NameToken, bot.IsProduction)
	if len(okxSignalId) == 0 {
		logger.Errorf("Error in creating a signal on the OKX exchange")
		u.Respond(w, u.Message(false, "Error in creating a signal on the OKX exchange"))
		return
	}
	time.Sleep(3 * time.Second)
	bot.OkxSignalId = okxSignalId
	instId := signal.NameToken + "-USDT-SWAP"
	strAmount := fmt.Sprintf("%2f", botRequest.InitialAmount)

	lever := strconv.FormatFloat(bot.Lever, 'f', 2, 64)
	bot.OkxBotId, err = OkxCreateSignalBot(user, okxSignalId, instId, lever, strAmount, bot.IsProduction)
	if err != nil {
		logger.Errorf("Error in OkxCreateSignalBot")
		u.Respond(w, u.Message(false, "Error in OkxCreateSignalBot"))
		return
	}

	bot.Status = models.Created
	bot.InitialAmount = botRequest.InitialAmount
	bot.CurrentAmount = botRequest.InitialAmount
	bot.StartTime = time.Now()

	resp := bot.Create(&signal)
	u.Respond(w, resp)
}

func getOkxSignalId(userId uint, strategyName string, tokenName string, isProduction bool) string {
	getSignalsRequest := new(model.GetSignalsRequest)
	getSignalsRequest.SignalSourceType = "1"
	signals := OkxGetSignals(userId, isProduction)

	if signals.Data != nil {
		for _, v := range signals.Data {
			signalChanId := v.SignalChanId
			bot, err := models.FindBotByByOkxSignalId(signalChanId)
			if err != nil {
				return ""
			}
			if bot.ID == 0 {
				logger.Info("getOkxSignalId: ", signalChanId)
				return signalChanId
			}
		}
	}
	nanoId := nanoid.New()[:4]
	okxSignalId := OkxCreateSignal(userId, strategyName+"("+nanoId+")", tokenName, isProduction)
	return okxSignalId
}

var DeleteBot = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)

	botDeleteRequest := &models.BotWithIdRequest{}

	err := json.NewDecoder(r.Body).Decode(botDeleteRequest)
	if err != nil {
		logger.Errorf("Error while decoding request body: %s", err.Error())
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	foundBot := models.Find(*botDeleteRequest)
	if foundBot.UserId == user {
		deal := models.FindByStatus(foundBot.ID, models.DealStarted)
		signalId := foundBot.SignalRefer
		signal, err := models.FindSignalById(signalId)
		if err == nil {
			if deal.ID > 0 {
				closeDeal(&deal, &foundBot, signal.NameToken)
			}
			OkxDeleteSignalBot(user, foundBot.OkxBotId, foundBot.IsProduction)
			foundBot.Status = models.Stopped
			foundBot.Update()
			foundBot.Delete()
			u.Respond(w, u.Message(true, "Success delete bot"))
		} else {
			logger.Errorf("Error while find signal: %s", err.Error())
			u.Respond(w, u.Message(false, "Wrong signal"))
		}
	} else {
		u.Respond(w, u.Message(false, "Wrong user"))
	}
}

var GetBots = func(w http.ResponseWriter, r *http.Request) {
	signalCode := r.URL.Query().Get("code")

	bots := models.GetBots(signalCode)
	resp := u.Message(true, "success")
	resp["bots"] = bots

	logger.Infoln("resp", resp)
	u.Respond(w, resp)
}

var GetAllOkxBots = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	bots := OkxGetAllActiveSignalBots(user, false)

	resp := u.Message(true, "success")
	resp["OKX bots"] = bots

	logger.Infoln("resp", resp)
	u.Respond(w, resp)
}

var GetAllOkxSignals = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	bots := OkxGetSignals(user, false)

	resp := u.Message(true, "success")
	resp["OKX signals"] = bots

	logger.Infoln("resp", resp)
	u.Respond(w, resp)
}
