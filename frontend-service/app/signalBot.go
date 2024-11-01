package app

import (
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx/futures"
)

func CreateSignal(api *futures.PrvApi, signalName string, signalDesc string) (signal *model.CreateSignalResponse, err error) {
	signalRequest := new(model.CreateSignalRequest)
	signalRequest.SignalChanName = signalName
	signalRequest.SignalChanDesc = signalDesc

	newSignal, data, err := api.Isolated.CreateSignal(*signalRequest)
	if err != nil {
		logger.Errorf("Error place order: %v, data: %v", err, string(data))
		return nil, err
	}
	signalId := newSignal.SignalChanId
	signalToken := newSignal.SignalChanToken
	logger.Info("signalId = ", signalId)
	logger.Info("signalToken = ", signalToken)

	return newSignal, nil
}

func CreateSignalBot(api *futures.PrvApi, signalChanId string, instIds string, lever string, investAmt string) (signal *model.CreateSignalBotResponse, err error) {
	createSignalBotRequest := new(model.CreateSignalBotRequest)
	createSignalBotRequest.SignalChanId = signalChanId
	createSignalBotRequest.Lever = lever
	createSignalBotRequest.InvestAmt = investAmt
	//createSignalBotRequest.IncludeAll = "true"
	createSignalBotRequest.InstIds = append(createSignalBotRequest.InstIds, instIds)
	createSignalBotRequest.SubOrdType = "2" //Sub order type 1：limit order 2：market order 9：tradingView signal

	entrySettingParamData := new(model.EntrySettingParamData)
	entrySettingParamData.EntryType = "1"
	entrySettingParamData.AllowMultipleEntry = "true"
	createSignalBotRequest.EntrySettingParam = *entrySettingParamData

	exitSettingParamData := new(model.ExitSettingParamData)
	createSignalBotRequest.ExitSettingParam = *exitSettingParamData

	newSignalBot, data, err := api.Isolated.CreateSignalBot(*createSignalBotRequest)
	if err != nil {
		logger.Errorf("Error creating signal bot: %v, data: %v", err, string(data))
		return nil, err
	}

	return newSignalBot, nil
}

func CancelSignalBot(api *futures.PrvApi, algoId string) (signal *model.CancelSignalBotResponse, err error) {
	cancelSignalBotRequest := new(model.CancelSignalBotRequest)
	cancelSignalBotRequest.AlgoId = algoId

	stopSignalBot, data, err := api.Isolated.CancelSignalBot(*cancelSignalBotRequest)
	if err != nil {
		logger.Errorf("Error CancelSignalBot: %v, data: %v", err, string(data))
		return nil, err
	}
	signalId := stopSignalBot.AlgoId
	logger.Info("signalId = ", signalId)

	return stopSignalBot, nil
}

func PlaceSubOrderSignalBot(api *futures.PrvApi, placeSubOrderSignalBotRequest *model.PlaceSubOrderSignalBotRequest) (signal *model.PlaceSubOrderSignalBotResponse, err error) {

	response, data, err := api.Isolated.PlaceSubOrderSignalBot(*placeSubOrderSignalBotRequest)
	if err != nil {
		logger.Errorf("Error PlaceSubOrderSignalBot: %v, data: %v", err, string(data))
		return nil, err
	}
	logger.Info("string(data) = ", string(data))

	return response, nil
}

func CancelSubOrderSignalBot(api *futures.PrvApi, cancelSubOrderSignalBot *model.CancelSubOrderSignalBotRequest) (signal *model.CancelSubOrderSignalBotResponse, err error) {

	response, data, err := api.Isolated.CancelSubOrderSignalBot(*cancelSubOrderSignalBot)
	if err != nil {
		logger.Errorf("Error CancelSubOrderSignalBotRequest: %v, data: %v", err, string(data))
		return nil, err
	}
	return response, nil
}

func ClosePositionSignalBot(api *futures.PrvApi, request *model.ClosePositionSignalBotRequest) (signal *model.ClosePositionSignalBotResponse, err error) {

	response, data, err := api.Isolated.ClosePositionSignalBot(*request)
	if err != nil {
		logger.Errorf("Error ClosePositionSignalBot: %v, data: %v", err, string(data))
		return nil, err
	}
	return response, nil
}

func GetSubOrderSignalBot(api *futures.PrvApi, request *model.GetSubOrdersSignalBotRequest) (signal *model.GetSubOrdersSignalBotResponse, err error) {

	response, data, err := api.Isolated.GetSubOrdersSignalBot(*request)
	if err != nil {
		logger.Errorf("Error GetSubOrderSignalBot: %v, data: %v", err, string(data))
		return nil, err
	}

	return response, nil
}

func GetSignals(api *futures.PrvApi, request *model.GetSignalsRequest) (signal *model.GetSignalsResponse, err error) {

	response, data, err := api.Isolated.GetSignals(*request)
	if err != nil {
		logger.Errorf("Error GetSubOrderSignalBot: %v, data: %v", err, string(data))
		return nil, err
	}

	return response, nil
}

func GetTicker(api *futures.PrvApi, symbol string) (signal *model.Ticker, err error) {

	response, data, err := api.Isolated.GetTicker(symbol)
	if err != nil {
		logger.Errorf("Error GetTicker: %v, data: %v", err, string(data))
		return nil, err
	}

	return response, nil
}

func GetActiveSignalBot(api *futures.PrvApi, algoId string) (signal *model.GetActiveSignalBotResponse, err error) {

	request := new(model.GetActiveSignalBotRequest)
	request.AlgoOrdType = "contract"
	response, data, err := api.Isolated.GetActiveSignalBot(*request)
	if err != nil {
		logger.Errorf("Error GetActiveSignalBot: %v, data: %v", err, string(data))
		return nil, err
	}

	return response, nil
}

func GetSignalBot(api *futures.PrvApi, algoId string) (signal *model.GetSignalBotResponse, err error) {

	request := new(model.GetSignalBotRequest)
	request.AlgoId = algoId
	request.AlgoOrdType = "contract"
	response, data, err := api.Isolated.GetSignalBot(*request)
	if err != nil {
		logger.Errorf("Error GetSignalBot: %v, data: %v", err, string(data))
		return nil, err
	}
	return response, nil
}

func GetAllActiveSignalBots(api *futures.PrvApi) (signal *model.GetActiveSignalBotResponse, err error) {

	request := new(model.GetActiveSignalBotRequest)
	request.AlgoOrdType = "contract"
	response, data, err := api.Isolated.GetActiveSignalBot(*request)
	if err != nil {
		logger.Errorf("Error GetActiveSignalBot: %v, data: %v", err, string(data))
		return nil, err
	}
	return response, nil
}
