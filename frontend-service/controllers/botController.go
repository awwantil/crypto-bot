package controllers

import (
	"encoding/json"
	"net/http"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
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

	resp := bot.Create(botRequest.CodeSignalId, botRequest.InitialAmount, botRequest.PosSide)
	u.Respond(w, resp)
}

var DeleteBot = func(w http.ResponseWriter, r *http.Request) {

	botDeleteRequest := &models.BotWithIdRequest{}

	err := json.NewDecoder(r.Body).Decode(botDeleteRequest)
	if err != nil {
		logger.Errorf("Error while decoding request body: %s", err.Error())
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	foundBot := models.Find(*botDeleteRequest)
	deal := models.FindByStatus(foundBot.ID, models.DealStarted)
	signalId := foundBot.SignalRefer
	signal, err := models.FindSignalById(signalId)
	if deal.ID > 0 {
		closeDeal(&deal, &foundBot, signal.NameToken)
	}
	foundBot.Delete()

	u.Respond(w, u.Message(true, "Success delete bot"))
}

var GetBots = func(w http.ResponseWriter, r *http.Request) {

	signalCode := r.URL.Query().Get("code")

	bots := models.GetBots(signalCode)
	resp := u.Message(true, "success")
	resp["bots"] = bots

	logger.Infoln("resp", resp)
	u.Respond(w, resp)
}
