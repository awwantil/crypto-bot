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
