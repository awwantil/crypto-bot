package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"okx-bot/exchange/logger"
	"okx-bot/exchange/model"
	"okx-bot/exchange/util"
)

func (prv *Prv) CreateSignal(req model.CreateSignalRequest, opt ...model.OptionParameter) (*model.CreateSignalResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.PostCreateSignalUri)

	params := url.Values{}
	params.Set("signalChanName", req.SignalChanName)
	params.Set("signalChanDesc", req.SignalChanDesc)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		logger.Errorf("[CreateSignal] err=%s, response=%s", err.Error(), string(data))
		return &model.CreateSignalResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	details, err := prv.UnmarshalOpts.PostCreateSignalUnmarshaler(data)

	return details, responseBody, err
}

func (prv *Prv) CreateSignalBot(req model.CreateSignalBotRequest, opt ...model.OptionParameter) (*model.CreateSignalBotResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.PostCreateSignalBotUri)

	jsonStr, _ := json.Marshal(req)

	data, responseBody, err := prv.DoAuthPostRequestWithParam(http.MethodPost, reqUrl, jsonStr, nil)
	if err != nil {
		logger.Errorf("[CreateSignalBot] err=%s, response=%s", err.Error(), string(data))
		return &model.CreateSignalBotResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	details, err := prv.UnmarshalOpts.PostCreateSignalBotUnmarshaler(data)

	return details, responseBody, err
}

func (prv *Prv) CancelSignalBot(req model.CancelSignalBotRequest, opt ...model.OptionParameter) (*model.CancelSignalBotResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.PostCancelSignalBotUri)
	var reqArray = []model.CancelSignalBotRequest{req}

	jsonStr, _ := json.Marshal(reqArray)
	data, responseBody, err := prv.DoAuthPostRequestWithParam(http.MethodPost, reqUrl, jsonStr, nil)
	if err != nil {
		logger.Errorf("[PostCancelSignalBotUri] err=%s, response=%s", err.Error(), string(data))
		return &model.CancelSignalBotResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	details, err := prv.UnmarshalOpts.PostCancelSignalBotUnmarshaler(data)

	return details, responseBody, err
}

func (prv *Prv) PlaceSubOrderSignalBot(req model.PlaceSubOrderSignalBotRequest, opt ...model.OptionParameter) (*model.PlaceSubOrderSignalBotResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.PostPlaceSubOrderSignalBotUri)

	params := url.Values{}
	params.Set("instId", req.InstId)
	params.Set("algoId", req.AlgoId)
	params.Set("side", req.Side)
	params.Set("ordType", req.OrdType)
	params.Set("sz", req.Sz)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		logger.Errorf("[PostPlaceSubOrderSignalBotUri] err=%s, response=%s", err.Error(), string(data))
		return &model.PlaceSubOrderSignalBotResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	details, err := prv.UnmarshalOpts.PostPlaceSubOrderSignalBotUnmarshaler(data)

	return details, responseBody, err
}

func (prv *Prv) CancelSubOrderSignalBot(req model.CancelSubOrderSignalBotRequest, opt ...model.OptionParameter) (*model.CancelSubOrderSignalBotResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.PostCancelSubOrderSignalBotUri)

	params := url.Values{}
	params.Set("instId", req.InstId)
	params.Set("algoId", req.AlgoId)
	params.Set("signalOrdId", req.SignalOrdId)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		logger.Errorf("[PostCancelSubOrderSignalBotUri] err=%s, response=%s", err.Error(), string(data))
		return &model.CancelSubOrderSignalBotResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data: ", string(data))

	details, err := prv.UnmarshalOpts.PostCancelSubOrderSignalBotUnmarshaler(data)

	return details, responseBody, err
}

func (prv *Prv) ClosePositionSignalBot(req model.ClosePositionSignalBotRequest, opt ...model.OptionParameter) (*model.ClosePositionSignalBotResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.PostClosePositionSignalBotUri)

	params := url.Values{}
	params.Set("instId", req.InstId)
	params.Set("algoId", req.AlgoId)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		logger.Errorf("[ClosePositionSignalBotRequest] err=%s, response=%s", err.Error(), string(data))
		return &model.ClosePositionSignalBotResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	details, err := prv.UnmarshalOpts.PostClosePositionSignalBotUnmarshaler(data)

	return details, responseBody, err
}

func (prv *Prv) GetSubOrdersSignalBot(req model.GetSubOrdersSignalBotRequest, opt ...model.OptionParameter) (*model.GetSubOrdersSignalBotResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetSubOrdersSignalBotUri)

	params := url.Values{}
	params.Set("algoId", req.AlgoId)
	params.Set("algoOrdType", req.AlgoOrdType)
	params.Set("signalOrdId", req.SignalOrdId)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	logger.Info("data for GetSubOrdersSignalBot: ", string(data))
	logger.Info("responseBody for GetSubOrdersSignalBot: ", string(responseBody))
	if err != nil {
		logger.Errorf("[GetSubOrdersSignalBotRequest] err=%s, response=%s", err.Error(), string(data))
		return &model.GetSubOrdersSignalBotResponse{}, responseBody, err
	}

	details, err := prv.UnmarshalOpts.GetSubOrdersSignalBotUnmarshaler(data)

	return details, responseBody, err
}