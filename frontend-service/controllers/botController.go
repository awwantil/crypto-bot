package controllers

import (
	"encoding/json"
	"net/http"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
)

var CreateBot = func(w http.ResponseWriter, r *http.Request) {

	botRequest := &models.BotCreateRequest{}

	err := json.NewDecoder(r.Body).Decode(botRequest)
	if err != nil {
		logger.Errorf("Error while decoding request body: %s", err.Error())
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	bot := new(models.Bot)

	resp := bot.Create(botRequest.CodeSignalId, botRequest.InitialAmount)
	u.Respond(w, resp)
}

var DeleteBot = func(w http.ResponseWriter, r *http.Request) {

	botDeleteRequest := &models.BotDeleteRequest{}

	err := json.NewDecoder(r.Body).Decode(botDeleteRequest)
	if err != nil {
		logger.Errorf("Error while decoding request body: %s", err.Error())
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	var foundBot = models.Bot{}
	db.Where("id = ?", botDeleteRequest.).First(&foundBot)

	logger.Infoln("foundBot", foundBot)

	resp := bot.Delete(botRequest.CodeSignalId, botRequest.InitialAmount)
	u.Respond(w, resp)
}

var GetBots = func(w http.ResponseWriter, r *http.Request) {

	signalCode := r.URL.Query().Get("code")

	bots := models.GetBots(signalCode)
	resp := u.Message(true, "success")
	resp["bots"] = bots

	logger.Infoln("resp", resp)
	u.Respond(w, resp)
}
