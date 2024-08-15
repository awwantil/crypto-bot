package common

import (
	"github.com/buger/jsonparser"
	"okx-bot/exchange/model"
)

func (un *RespUnmarshaler) UnmarshalCreateSignal(data []byte) (*model.CreateSignalResponse, error) {
	var details = new(model.CreateSignalResponse)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "signalChanId":
				details.SignalChanId = detailsStr
			case "signalChanToken":
				details.SignalChanToken = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return details, err
}
