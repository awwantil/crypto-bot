package common

import (
	"github.com/buger/jsonparser"
	"okx-bot/exchange/logger"
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

func (un *RespUnmarshaler) UnmarshalCreateSignalBot(data []byte) (*model.CreateSignalBotResponse, error) {
	var details = new(model.CreateSignalBotResponse)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "algoId":
				details.AlgoId = detailsStr
			case "algoClOrdId":
				details.AlgoClOrdId = detailsStr
			case "sCode":
				details.SCode = detailsStr
			case "sMsg":
				details.SMsg = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return details, err
}

func (un *RespUnmarshaler) UnmarshalCancelSignalBot(data []byte) (*model.CancelSignalBotResponse, error) {
	var details = new(model.CancelSignalBotResponse)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "algoId":
				details.AlgoId = detailsStr
			case "sCode":
				details.SCode = detailsStr
			case "sMsg":
				details.SMsg = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return details, err
}

func (un *RespUnmarshaler) UnmarshalPlaceSubOrderSignalBot(data []byte) (*model.PlaceSubOrderSignalBotResponse, error) {
	var details = new(model.PlaceSubOrderSignalBotResponse)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "code":
				details.Code = detailsStr
			case "msg":
				details.Msg = detailsStr
			case "data":
				details.Data = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return details, err
}

func (un *RespUnmarshaler) UnmarshalCancelSubOrderSignalBot(data []byte) (*model.CancelSubOrderSignalBotResponse, error) {
	var details = new(model.CancelSubOrderSignalBotResponse)
	var detailsData = new(model.CancelSubOrderData)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "code":
				details.Code = detailsStr
			case "msg":
				details.Msg = detailsStr
			case "data":
				logger.Infof("data 1 - !!!!!!!!!!!! %v", detailsStr)
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						logger.Infof("data 2 - !!!!!!!!!!!! %v", detailsStr)
						switch string(key) {
						case "signalOrdId":
							detailsData.SignalOrdId = detailsStr
						case "sCode":
							detailsData.SCode = detailsStr
						case "sMsg":
							detailsData.SMsg = detailsStr
						}
						return err
					})
					if err != nil {
						return
					}
					details.Data = append(details.Data, *detailsData)
				})
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return details, err
}

func (un *RespUnmarshaler) UnmarshalClosePositionSignalBot(data []byte) (*model.ClosePositionSignalBotResponse, error) {
	var details = new(model.ClosePositionSignalBotResponse)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "algoId":
				details.AlgoId = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return details, err
}

func (un *RespUnmarshaler) UnmarshalGetSubOrdersSignalBot(data []byte) (*model.GetSubOrdersSignalBotResponse, error) {
	var details = new(model.GetSubOrdersSignalBotResponse)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "algoId":
				details.AlgoId = detailsStr
			case "instId":
				details.InstId = detailsStr
			case "instType":
				details.InstType = detailsStr
			case "uTime":
				details.UTime = detailsStr
			case "cTime":
				details.CTime = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return details, err
}
