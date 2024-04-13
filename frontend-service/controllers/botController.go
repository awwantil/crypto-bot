package controllers

import (
	"encoding/json"
	"net/http"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
)

var CreateBot = func(w http.ResponseWriter, r *http.Request) {

	botRequest := &models.BotRequest{}

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
