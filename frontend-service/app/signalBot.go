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
	createSignalBotRequest.InstIds = append(createSignalBotRequest.InstIds, instIds)
	createSignalBotRequest.SubOrdType = "9" //Sub order type 1：limit order 2：market order 9：tradingView signal

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
	algoId := newSignalBot.AlgoId
	sCode := newSignalBot.SCode
	sMsg := newSignalBot.SMsg
	logger.Info("algoId = ", algoId)
	logger.Info("sCode = ", sCode)
	logger.Info("sMsg = ", sMsg)

	return newSignalBot, nil
}
