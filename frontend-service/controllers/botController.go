package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	bot := new(models.Bot)
	bot.UserId = user

	var signal = models.Signal{Code: botRequest.CodeSignalId}
	models.GetDB().Where("code = ?", botRequest.CodeSignalId).First(&signal)

	logger.Infoln("signal", signal)

	okxSignalId := OkxCreateSignal(bot.UserId, signal.StrategyName, signal.NameToken)
	bot.OkxSignalId = okxSignalId
	instId := signal.NameToken + "-USDT-SWAP"
	strAmount := fmt.Sprintf("%2f", botRequest.InitialAmount)
	bot.OkxBotId = OkxCreateSignalBot(bot.UserId, okxSignalId, instId, "3", strAmount)
	bot.Status = models.Created
	bot.InitialAmount = botRequest.InitialAmount
	bot.CurrentAmount = botRequest.InitialAmount
	bot.PosSide = botRequest.PosSide
	bot.StartTime = time.Now()

	resp := bot.Create(&signal)
	u.Respond(w, resp)
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
