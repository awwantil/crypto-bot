package controllers

import (
	"encoding/json"
	"net/http"
	"okx-bot/frontend-service/app"
	"okx-bot/frontend-service/models"
	u "okx-bot/frontend-service/utils"
)

var CreateOkxApi = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	okxApi := &models.OKxApi{}

	err := json.NewDecoder(r.Body).Decode(okxApi)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	okxApi.UserId = user
	resp := okxApi.Create()
	u.Respond(w, resp)
}

var GetOkxApiFor = func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user").(uint)
	data, err := models.GetUserApi(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	if err != nil {
		resp = u.Message(false, "error")
	}
	u.Respond(w, resp)
}

func OkxCreateSignal(userId uint, signalName string, signalDesc string) string {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return ""
	}
	newSignal, err := app.CreateSignal(api, signalName, signalDesc)
	if err != nil {
		return ""
	}
	return newSignal.SignalChanId
}

func OkxCreateSignalBot(userId uint, signalChanId string, instIds string, lever string, investAmt string) string {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return ""
	}
	newSignalBot, err := app.CreateSignalBot(api, signalChanId, instIds, lever, investAmt)
	if err != nil {
		return ""
	}
	return newSignalBot.AlgoId
}

func OkxDeleteSignalBot(userId uint, signalChanId string) string {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return ""
	}
	cancelSignalBot, err := app.CancelSignalBot(api, signalChanId)
	if err != nil {
		return ""
	}
	return cancelSignalBot.AlgoId
}
