package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aidarkhanov/nanoid"
	"net/http"
	"okx-bot/exchange/model"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
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

	bot := new(models.Bot)
	bot.UserId = user

	var signal = models.Signal{Code: botRequest.CodeSignalId}
	models.GetDB().Where("code = ?", botRequest.CodeSignalId).First(&signal)

	logger.Infoln("signal", signal)

	okxSignalId := getOkxSignalId(user, signal.StrategyName, signal.NameToken)
	if len(okxSignalId) == 0 {
		logger.Errorf("Error in creating a signal on the OKX exchange")
		u.Respond(w, u.Message(false, "Error in creating a signal on the OKX exchange"))
		return
	}
	time.Sleep(3 * time.Second)
	bot.OkxSignalId = okxSignalId
	instId := signal.NameToken + "-USDT-SWAP"
	strAmount := fmt.Sprintf("%2f", botRequest.InitialAmount)

	bot.OkxBotId, err = OkxCreateSignalBot(user, okxSignalId, instId, "3", strAmount)
	if err != nil {
		logger.Errorf("Error in OkxCreateSignalBot")
		u.Respond(w, u.Message(false, "Error in OkxCreateSignalBot"))
		return
	}

	bot.Status = models.Created
	bot.InitialAmount = botRequest.InitialAmount
	bot.CurrentAmount = botRequest.InitialAmount
	bot.PosSide = botRequest.PosSide
	bot.StartTime = time.Now()

	resp := bot.Create(&signal)
	u.Respond(w, resp)
}

func getOkxSignalId(userId uint, strategyName string, tokenName string) string {
	getSignalsRequest := new(model.GetSignalsRequest)
	getSignalsRequest.SignalSourceType = "1"
	signals := OkxGetSignals(userId)
	logger.Info("signals", signals.Data[0].SignalChanId)

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
	nanoId := nanoid.New()[:4]
	okxSignalId := OkxCreateSignal(userId, strategyName+"("+nanoId+")", tokenName)
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
			OkxDeleteSignalBot(user, foundBot.OkxBotId)
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
	user := r.Context().Value("user").(uint)
	ticker := OkxGetTicker(user, "SOL")
	logger.Info("ticker", ticker)
	signalCode := r.URL.Query().Get("code")

	bots := models.GetBots(signalCode)
	resp := u.Message(true, "success")
	resp["bots"] = bots

	logger.Infoln("resp", resp)
	u.Respond(w, resp)
}
