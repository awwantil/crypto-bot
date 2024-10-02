package controllers

import (
	"encoding/json"
	"net/http"
	"okx-bot/exchange/model"
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

func OkxCreateSignalBot(userId uint, signalChanId string, instIds string, lever string, investAmt string) (string, error) {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return "", err
	}
	newSignalBot, err := app.CreateSignalBot(api, signalChanId, instIds, lever, investAmt)
	if err != nil {
		logger.Errorf("Error in OkxCreateSignalBot: %v", err)
		return "", err
	}
	return newSignalBot.AlgoId, nil
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

func OkxPlaceSubOrder(userId uint, instId string, algoId string, sz string) (string, error) {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return "", err
	}
	placeSubOrderSignalBotRequest := new(model.PlaceSubOrderSignalBotRequest)
	placeSubOrderSignalBotRequest.InstId = instId
	placeSubOrderSignalBotRequest.AlgoId = algoId
	placeSubOrderSignalBotRequest.Side = "buy"
	placeSubOrderSignalBotRequest.OrdType = "market"
	placeSubOrderSignalBotRequest.Sz = sz

	response, err := app.PlaceSubOrderSignalBot(api, placeSubOrderSignalBotRequest)
	if err != nil {
		return "", err
	}
	return response.Code, nil
}

func OkxGetSubOrderSignalBot(userId uint, algoId string) string {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return ""
	}
	getSubOrdersSignalBotRequest := new(model.GetSubOrdersSignalBotRequest)
	getSubOrdersSignalBotRequest.AlgoId = algoId
	getSubOrdersSignalBotRequest.AlgoOrdType = "contract"
	getSubOrdersSignalBotRequest.State = "filled"

	details, err := app.GetSubOrderSignalBot(api, getSubOrdersSignalBotRequest)
	if err != nil {
		return ""
	}
	logger.Infof("details OrdId: %v", details.OrdId)
	return details.OrdId
}

func OkxGetSignals(userId uint) *model.GetSignalsResponse {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return nil
	}
	getSignalsRequest := new(model.GetSignalsRequest)
	getSignalsRequest.SignalSourceType = "1"

	details, err := app.GetSignals(api, getSignalsRequest)
	if err != nil {
		return nil
	}
	logger.Infof("details OrdId: %v", details.Data)
	return details
}

func OkxGetTicker(userId uint, symbol string) *model.Ticker {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return nil
	}
	details, err := app.GetTicker(api, symbol)
	if err != nil {
		return nil
	}
	logger.Infof("details OrdId: %v", details.Last)
	return details
}

func OkxGetActiveSignalBot(userId uint, algoId string) *model.GetActiveSignalBotResponseData {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return nil
	}
	details, err := app.GetActiveSignalBot(api, algoId)
	if err != nil {
		return nil
	}
	var result = new(model.GetActiveSignalBotResponseData)
	for _, v := range details.Bots {
		if v.AlgoId == algoId {
			logger.Info("found bot with algoId: ", algoId)
			result = &v
		}
	}
	return result
}

func OkxGetSignalBot(userId uint, algoId string) *model.GetSignalBotResponseData {
	api, err := app.GetOkxApi(userId)
	if err != nil {
		logger.Errorf("Error in GetOkxApi: %v", err)
		return nil
	}
	details, err := app.GetSignalBot(api, algoId)
	if err != nil {
		return nil
	}
	var result = new(model.GetSignalBotResponseData)
	for _, v := range details.Bots {
		if v.AlgoId == algoId {
			logger.Info("found bot with algoId: ", algoId)
			result = &v
		}
	}
	return result
}
