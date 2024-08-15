package signal_bot

import (
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx/common"
)

type PrvApi struct {
	*common.Prv
}

func (prv *PrvApi) CreateSignal(req model.CreateSignalRequest) (*model.CreateSignalResponse, []byte, error) {
	details, respBody, err := prv.Prv.CreateSignal(req)
	return details, respBody, err
}
