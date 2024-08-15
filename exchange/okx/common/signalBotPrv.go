package common

import (
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
