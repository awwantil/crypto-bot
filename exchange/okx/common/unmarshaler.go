package common

import (
	"encoding/json"
	"fmt"
	"okx-bot/exchange/logger"
	"okx-bot/exchange/model"
	"time"

	"github.com/buger/jsonparser"
	"github.com/spf13/cast"
)

type RespUnmarshaler struct {
}

func (un *RespUnmarshaler) UnmarshalDepth(data []byte) (*model.Depth, error) {
	var (
		dep model.Depth
		err error
	)

	err = jsonparser.ObjectEach(data[1:len(data)-1],
		func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			switch string(key) {
			case "ts":
				dep.UTime = time.UnixMilli(cast.ToInt64(string(value)))
			case "asks":
				items, _ := un.unmarshalDepthItem(value)
				dep.Asks = items
			case "bids":
				items, _ := un.unmarshalDepthItem(value)
				dep.Bids = items
			}
			return nil
		})

	return &dep, err
}

func (un *RespUnmarshaler) unmarshalDepthItem(data []byte) (model.DepthItems, error) {
	var items model.DepthItems
	_, err := jsonparser.ArrayEach(data, func(asksItemData []byte, dataType jsonparser.ValueType, offset int, err error) {
		item := model.DepthItem{}
		i := 0
		_, err = jsonparser.ArrayEach(asksItemData, func(itemVal []byte, dataType jsonparser.ValueType, offset int, err error) {
			valStr := string(itemVal)
			switch i {
			case 0:
				item.Price = cast.ToFloat64(valStr)
			case 1:
				item.Amount = cast.ToFloat64(valStr)
			}
			i += 1
		})
		items = append(items, item)
	})
	return items, err
}

func (un *RespUnmarshaler) UnmarshalTicker(data []byte) (*model.Ticker, error) {
	var tk = &model.Ticker{}

	var open float64
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, val []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(val)
			switch string(key) {
			case "last":
				tk.Last = cast.ToFloat64(valStr)
			case "askPx":
				tk.Sell = cast.ToFloat64(valStr)
			case "bidPx":
				tk.Buy = cast.ToFloat64(valStr)
			case "vol24h":
				tk.Vol = cast.ToFloat64(valStr)
			case "high24h":
				tk.High = cast.ToFloat64(valStr)
			case "low24h":
				tk.Low = cast.ToFloat64(valStr)
			case "ts":
				tk.Timestamp = cast.ToInt64(valStr)
			case "open24h":
				open = cast.ToFloat64(valStr)
			}
			return nil
		})
	})

	if err != nil {
		logger.Errorf("[UnmarshalTicker] %s", err.Error())
		return nil, err
	}

	tk.Percent = (tk.Last - open) / open * 100

	return tk, nil
}

func (un *RespUnmarshaler) UnmarshalGetKlineResponse(data []byte) ([]model.Kline, error) {
	var klines []model.Kline
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var (
			k model.Kline
			i int
		)
		_, err = jsonparser.ArrayEach(value, func(val []byte, dataType jsonparser.ValueType, offset int, err error) {
			valStr := string(val)
			switch i {
			case 0:
				k.Timestamp = cast.ToInt64(valStr)
			case 1:
				k.Open = cast.ToFloat64(valStr)
			case 2:
				k.High = cast.ToFloat64(valStr)
			case 3:
				k.Low = cast.ToFloat64(valStr)
			case 4:
				k.Close = cast.ToFloat64(valStr)
			case 5:
				k.Vol = cast.ToFloat64(valStr)
			}
			i += 1
		})
		klines = append(klines, k)
	})

	return klines, err
}

func (un *RespUnmarshaler) UnmarshalCreateOrderResponse(data []byte) (*model.Order, error) {
	var ord = new(model.Order)
	err := jsonparser.ObjectEach(data[1:len(data)-1], func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		valStr := string(value)
		switch string(key) {
		case "ordId":
			ord.Id = valStr
		case "clOrdId":
			ord.CId = valStr
		}
		return nil
	})
	return ord, err
}

func (un *RespUnmarshaler) UnmarshalSetLeverageResponse(data []byte) (*model.SetLeverageResponse, error) {
	var ord = new(model.SetLeverageResponse)
	err := jsonparser.ObjectEach(data[1:len(data)-1], func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		valStr := string(value)
		switch string(key) {
		case "lever":
			ord.Lever = valStr
		case "instId":
			ord.InstId = valStr
		}
		return nil
	})
	return ord, err
}

func (un *RespUnmarshaler) UnmarshalGetPendingOrdersResponse(data []byte) ([]model.OrderInfoResponse, error) {
	var (
		orders []model.OrderInfoResponse
		err    error
	)

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ord, err := un.UnmarshalGetOrderInfoResponse(value)
		if err != nil {
			return
		}
		orders = append(orders, *ord)
	})

	return orders, err
}

func (un *RespUnmarshaler) UnmarshalGetHistoryOrdersResponse(data []byte) ([]model.OrderInfoResponse, error) {
	return un.UnmarshalGetPendingOrdersResponse(data)
}

func (un *RespUnmarshaler) UnmarshalGetOrderInfoResponse(data []byte) (ord *model.OrderInfoResponse, err error) {
	var orderInfoResponse = new(model.OrderInfoResponse)
	var attachAlgoOrds = new(model.AttachAlgoOrds)

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "instType":
				orderInfoResponse.InstType = detailsStr
			case "instId":
				orderInfoResponse.InstId = detailsStr
			case "tgtCcy":
				orderInfoResponse.TgtCcy = detailsStr
			case "ccy":
				orderInfoResponse.Ccy = detailsStr
			case "ordId":
				orderInfoResponse.OrdId = detailsStr
			case "clOrdId":
				orderInfoResponse.ClOrdId = detailsStr
			case "tag":
				orderInfoResponse.Tag = detailsStr
			case "px":
				orderInfoResponse.Px = detailsStr
			case "pxUsd":
				orderInfoResponse.PxUsd = detailsStr
			case "pxVol":
				orderInfoResponse.PxVol = detailsStr
			case "pxType":
				orderInfoResponse.PxType = detailsStr
			case "sz":
				orderInfoResponse.Sz = detailsStr
			case "pnl":
				orderInfoResponse.Pnl = detailsStr
			case "ordType":
				orderInfoResponse.OrdType = detailsStr
			case "side":
				orderInfoResponse.Side = detailsStr
			case "posSide":
				orderInfoResponse.PosSide = detailsStr
			case "tdMode":
				orderInfoResponse.TdMode = detailsStr
			case "accFillSz":
				orderInfoResponse.AccFillSz = detailsStr
			case "fillPx":
				orderInfoResponse.FillPx = detailsStr
			case "tradeId":
				orderInfoResponse.TradeId = detailsStr
			case "fillSz":
				orderInfoResponse.FillSz = detailsStr
			case "fillTime":
				orderInfoResponse.FillTime = detailsStr
			case "avgPx":
				orderInfoResponse.AvgPx = detailsStr
			case "state":
				orderInfoResponse.State = detailsStr
			case "stpMode":
				orderInfoResponse.StpMode = detailsStr
			case "lever":
				orderInfoResponse.Lever = detailsStr
			case "attachAlgoClOrdId":
				orderInfoResponse.AttachAlgoClOrdId = detailsStr
			case "tpTriggerPx":
				orderInfoResponse.TpTriggerPx = detailsStr
			case "tpOrdPx":
				orderInfoResponse.TpOrdPx = detailsStr
			case "slTriggerPx":
				orderInfoResponse.SlTriggerPx = detailsStr
			case "slTriggerPxType":
				orderInfoResponse.SlTriggerPxType = detailsStr
			case "slOrdPx":
				orderInfoResponse.SlOrdPx = detailsStr
			case "cTime":
				orderInfoResponse.CTime = time.Unix(cast.ToInt64(detailsStr)/1000, 0).Local()
			case "uTime":
				orderInfoResponse.UTime = time.Unix(cast.ToInt64(detailsStr)/1000, 0).Local()
			case "attachAlgoOrds":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						switch string(key) {
						case "attachAlgoId":
							attachAlgoOrds.AttachAlgoId = detailsStr
						case "attachAlgoClOrdId":
							attachAlgoOrds.AttachAlgoClOrdId = detailsStr
						case "tpOrdKind":
							attachAlgoOrds.TtpOrdKind = detailsStr
						case "tpTriggerPx":
							attachAlgoOrds.TpTriggerPx = detailsStr
						case "tpTriggerPxType":
							attachAlgoOrds.TpTriggerPxType = detailsStr
						case "tpOrdPx":
							attachAlgoOrds.TpOrdPx = detailsStr
						case "slTriggerPx":
							attachAlgoOrds.SlTriggerPx = detailsStr
						case "slTriggerPxType":
							attachAlgoOrds.SlTriggerPxType = detailsStr
						case "slOrdPx":
							attachAlgoOrds.SlOrdPx = detailsStr
						case "sz":
							attachAlgoOrds.Sz = detailsStr
						case "amendPxOnTriggerType":
							attachAlgoOrds.AmendPxOnTriggerType = detailsStr
						}
						return err
					})
					if err != nil {
						return
					}
					orderInfoResponse.AttachAlgoOrds = append(orderInfoResponse.AttachAlgoOrds, *attachAlgoOrds)
				})
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return orderInfoResponse, err
}

func (un *RespUnmarshaler) UnmarshalGetAccountResponse(data []byte) (map[string]model.Account, error) {
	var accMap = make(map[string]model.Account, 2)

	_, err := jsonparser.ArrayEach(data[1:len(data)-1], func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var acc model.Account
		err = jsonparser.ObjectEach(value, func(key []byte, accData []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(accData)
			switch string(key) {
			case "ccy":
				acc.Coin = valStr
			case "availEq":
				acc.AvailableBalance = cast.ToFloat64(valStr)
			case "eq":
				acc.Balance = cast.ToFloat64(valStr)
			case "frozenBal":
				acc.FrozenBalance = cast.ToFloat64(valStr)
			}
			return err
		})

		if err != nil {
			return
		}

		accMap[acc.Coin] = acc
	}, "details")

	return accMap, err
}

func (un *RespUnmarshaler) UnmarshalGetFuturesAccountResponse(data []byte) (map[string]model.FuturesAccount, error) {
	var accMap = make(map[string]model.FuturesAccount, 2)

	_, err := jsonparser.ArrayEach(data[1:len(data)-1], func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var acc model.FuturesAccount
		err = jsonparser.ObjectEach(value, func(key []byte, accData []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(accData)
			switch string(key) {
			case "ccy":
				acc.Coin = valStr
			case "availEq":
				acc.AvailEq = cast.ToFloat64(valStr)
			case "eq":
				acc.Eq = cast.ToFloat64(valStr)
			case "frozenBal":
				acc.FrozenBal = cast.ToFloat64(valStr)
			case "upl":
				acc.Upl = cast.ToFloat64(valStr)
			case "mgnRatio":
				acc.MgnRatio = cast.ToFloat64(valStr)
			}
			return err
		})

		if err != nil {
			return
		}

		accMap[acc.Coin] = acc
	}, "details")

	return accMap, err
}

func (un *RespUnmarshaler) UnmarshalCancelOrderResponse(data []byte) (model.CancelOrderResponse, error) {
	var cancelOrderResponse = new(model.CancelOrderResponse)
	var cancelOrderResponseData = new(model.OrderResponseData)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "code":
				cancelOrderResponse.Code = detailsStr
			case "msg":
				cancelOrderResponse.Msg = detailsStr
			case "inTime":
				cancelOrderResponse.InTime = detailsStr
			case "outTime":
				cancelOrderResponse.OutTime = detailsStr
			case "data":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						switch string(key) {
						case "ordId":
							cancelOrderResponseData.OrdId = detailsStr
						case "clOrdId":
							cancelOrderResponseData.ClOrdId = detailsStr
						case "sCode":
							cancelOrderResponseData.SCode = detailsStr
						case "sMsg":
							cancelOrderResponseData.SMsg = detailsStr
						}
						return err
					})
					if err != nil {
						return
					}
					cancelOrderResponse.Data = append(cancelOrderResponse.Data, *cancelOrderResponseData)
				})
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return *cancelOrderResponse, err
}

func (un *RespUnmarshaler) UnmarshalClosePositionsResponse(data []byte) (model.ClosePositionResponse, error) {
	var closePositionResponse = new(model.ClosePositionResponse)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "instId":
				closePositionResponse.InstId = detailsStr
			case "posSide":
				closePositionResponse.PosSide = detailsStr
			case "clOrdId":
				closePositionResponse.ClOrdId = detailsStr
			case "tag":
				closePositionResponse.Tag = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return *closePositionResponse, err
}

func (un *RespUnmarshaler) UnmarshalGetPositionsResponse(data []byte) ([]model.FuturesPosition, error) {
	var (
		positions []model.FuturesPosition
		err       error
	)

	_, err = jsonparser.ArrayEach(data, func(posData []byte, dataType jsonparser.ValueType, offset int, err error) {
		var pos model.FuturesPosition
		err = jsonparser.ObjectEach(posData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(value)
			switch string(key) {
			case "availPos":
				pos.AvailQty = cast.ToFloat64(valStr)
			case "avgPx":
				pos.AvgPx = cast.ToFloat64(valStr)
			case "pos":
				pos.Qty = cast.ToFloat64(valStr)
			case "posSide":
				if valStr == "long" {
					pos.PosSide = model.Futures_OpenBuy
				}
				if valStr == "short" {
					pos.PosSide = model.Futures_OpenSell
				}
			case "upl":
				pos.Upl = cast.ToFloat64(valStr)
			case "uplRatio":
				pos.UplRatio = cast.ToFloat64(valStr)
			case "lever":
				pos.Lever = cast.ToFloat64(valStr)
			}
			return nil
		})
		positions = append(positions, pos)
	})

	return positions, err
}

func (un *RespUnmarshaler) UnmarshalGetPositionsHisotoryResponse(data []byte) ([]model.FuturesPositionHistory, error) {
	var (
		positionsHistory []model.FuturesPositionHistory
		err              error
	)

	_, err = jsonparser.ArrayEach(data, func(posData []byte, dataType jsonparser.ValueType, offset int, err error) {
		var posHistory model.FuturesPositionHistory
		_ = jsonparser.ObjectEach(posData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(value)
			switch string(key) {
			case "instId":
				posHistory.InstId = cast.ToString(valStr)
			case "direction":
				posHistory.Direction = cast.ToString(valStr)
			case "lever":
				posHistory.Lever = cast.ToFloat32(valStr)
			case "type":
				posHistory.Type = cast.ToString(valStr)
			case "cTime":
				posHistory.CTime = time.Unix(cast.ToInt64(valStr)/1000, 0).Local()
			case "uTime":
				posHistory.UTime = time.Unix(cast.ToInt64(valStr)/1000, 0).Local()
			case "openAvgPx":
				posHistory.OpenAvgPx = cast.ToFloat64(valStr)
			case "closeAvgPx":
				posHistory.CloseAvgPx = cast.ToFloat64(valStr)
			case "pnl":
				posHistory.Pnl = cast.ToFloat64(valStr)
			case "realizedPnl":
				posHistory.RealizedPnl = cast.ToFloat64(valStr)
			}
			return nil
		})
		positionsHistory = append(positionsHistory, posHistory)
	})
	return positionsHistory, err
}

func (un *RespUnmarshaler) UnmarshalGetAccountBalanceResponse(data []byte) (model.BalanceResponse, error) {
	logger.Info(string(data))

	var balanceResponse = new(model.BalanceResponse)
	var detailsResponseData = new(model.BalanceDetails)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "uTime":
				balanceResponse.UTime = detailsStr
			case "totalEq":
				balanceResponse.TotalEq = detailsStr
			case "isoEq":
				balanceResponse.IsoEq = detailsStr
			case "adjEq":
				balanceResponse.AdjEq = detailsStr
			case "ordFroz":
				balanceResponse.OrdFroz = detailsStr
			case "imr":
				balanceResponse.Imr = detailsStr
			case "mmr":
				balanceResponse.Mmr = detailsStr
			case "borrowFroz":
				balanceResponse.BorrowFroz = detailsStr
			case "mgnRatio":
				balanceResponse.MgnRatio = detailsStr
			case "notionalUsd":
				balanceResponse.NotionalUsd = detailsStr
			case "upl":
				balanceResponse.Upl = detailsStr
			case "details":
				_, err = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						switch string(key) {
						case "ccy":
							detailsResponseData.Ccy = detailsStr
						case "eq":
							detailsResponseData.Eq = detailsStr
						case "cashBal":
							detailsResponseData.CashBal = detailsStr
						case "uTime":
							detailsResponseData.UTime = detailsStr
						case "isoEq":
							detailsResponseData.IsoEq = detailsStr
						case "availEq":
							detailsResponseData.AvailEq = detailsStr
						case "disEq":
							detailsResponseData.DisEq = detailsStr
						case "fixedBal":
							detailsResponseData.FixedBal = detailsStr
						case "availBal":
							detailsResponseData.AvailBal = detailsStr
						case "frozenBal":
							detailsResponseData.FrozenBal = detailsStr
						case "ordFrozen":
							detailsResponseData.OrdFrozen = detailsStr
						case "liab":
							detailsResponseData.Liab = detailsStr
						case "upl":
							detailsResponseData.Upl = detailsStr
						case "uplLiab":
							detailsResponseData.UplLiab = detailsStr
						case "crossLiab":
							detailsResponseData.CrossLiab = detailsStr
						case "rewardBal":
							detailsResponseData.RewardBal = detailsStr
						case "isoLiab":
							detailsResponseData.IsoLiab = detailsStr
						case "mgnRatio":
							detailsResponseData.MgnRatio = detailsStr
						case "interest":
							detailsResponseData.Interest = detailsStr
						case "twap":
							detailsResponseData.Twap = detailsStr
						case "maxLoan":
							detailsResponseData.MaxLoan = detailsStr
						case "eqUsd":
							detailsResponseData.EqUsd = detailsStr
						case "borrowFroz":
							detailsResponseData.BorrowFroz = detailsStr
						case "notionalLever":
							detailsResponseData.NotionalLever = detailsStr
						case "stgyEq":
							detailsResponseData.StgyEq = detailsStr
						case "isoUpl":
							detailsResponseData.IsoUpl = detailsStr
						case "spotInUseAmt":
							detailsResponseData.SpotInUseAmt = detailsStr
						case "clSpotInUseAmt":
							detailsResponseData.ClSpotInUseAmt = detailsStr
						case "maxSpotInUseAmt":
							detailsResponseData.MaxSpotInUseAmt = detailsStr
						case "spotIsoBal":
							detailsResponseData.SpotIsoBal = detailsStr
						case "imr":
							detailsResponseData.Imr = detailsStr
						case "mmr":
							detailsResponseData.Mmr = detailsStr
						case "smtSyncEq":
							detailsResponseData.SmtSyncEq = detailsStr
						}
						return err
					})
					if err != nil {
						logger.Error(err)
						return
					}
					balanceResponse.Details = append(balanceResponse.Details, *detailsResponseData)
				})
			}
			return err
		})
		if err != nil {
			logger.Error(err)
			return
		}
	})
	if err != nil {
		logger.Error(err)
		return *balanceResponse, err
	}

	return *balanceResponse, nil
}

func (un *RespUnmarshaler) UnmarshalGetExchangeInfoResponse(data []byte) (map[string]model.CurrencyPair, error) {
	var (
		err             error
		currencyPairMap = make(map[string]model.CurrencyPair, 20)
	)

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var (
			currencyPair model.CurrencyPair
			instTy       string
			ctValCcy     string
			settleCcy    string
		)

		err = jsonparser.ObjectEach(value, func(key []byte, val []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(val)
			switch string(key) {
			case "instType":
				instTy = valStr
			case "instId":
				currencyPair.Symbol = valStr
			case "minSz":
				currencyPair.MinQty = cast.ToFloat64(valStr)
			case "tickSz":
				currencyPair.PricePrecision = AdaptQtyOrPricePrecision(valStr)
			case "lotSz":
				currencyPair.QtyPrecision = AdaptQtyOrPricePrecision(valStr)
			case "baseCcy":
				currencyPair.BaseSymbol = valStr
			case "quoteCcy":
				currencyPair.QuoteSymbol = valStr
			case "ctValCcy":
				ctValCcy = valStr
				currencyPair.ContractValCurrency = valStr
			case "ctVal":
				currencyPair.ContractVal = cast.ToFloat64(valStr)
			case "settleCcy":
				settleCcy = valStr
				currencyPair.SettlementCurrency = valStr
			case "alias":
				currencyPair.ContractAlias = valStr
			case "expTime":
				currencyPair.ContractDeliveryDate = cast.ToInt64(valStr)
			}
			return nil
		})

		if instTy == "SWAP" {
			currencyPair.BaseSymbol = ctValCcy
			currencyPair.QuoteSymbol = settleCcy
		}

		//adapt
		if instTy == "FUTURES" {
			currencyPair.BaseSymbol = settleCcy
			currencyPair.QuoteSymbol = ctValCcy
		}

		k := fmt.Sprintf("%s%s%s", currencyPair.BaseSymbol, currencyPair.QuoteSymbol, currencyPair.ContractAlias)
		currencyPairMap[k] = currencyPair
	})

	return currencyPairMap, err
}

func (un *RespUnmarshaler) UnmarshalResponse(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}

func (un *RespUnmarshaler) UnmarshalGetComputeMinInvestmentResponse(data []byte) (model.ComputeMinInvestmentResponse, error) {
	var minInvestment = new(model.ComputeMinInvestmentResponse)
	var investData = new(model.InvestmentData)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			investmentDataStr := string(respData)
			switch string(key) {
			case "minInvestmentData":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						newInvestmentDataStr := string(respData)
						switch string(key) {
						case "amt":
							investData.Amt = newInvestmentDataStr
						case "ccy":
							investData.Ccy = newInvestmentDataStr
						}
						return err
					})
					if err != nil {
						return
					}
				})
			case "singleAmt":
				minInvestment.SingleAmt = investmentDataStr
			}
			return err
		})

		if err != nil {
			return
		}
	})
	minInvestment.InvestmentData = append(minInvestment.InvestmentData, *investData)

	return *minInvestment, err
}

func (un *RespUnmarshaler) UnmarshalGetAlgoOrderDetailsResponse(data []byte) (model.GridAlgoOrderDetailsResponse, error) {
	var details = new(model.GridAlgoOrderDetailsResponse)
	var rebateTrans = new(model.RebateTrans)
	var triggerParams = new(model.TriggerParams)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "rebateTrans":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						switch string(key) {
						case "rebate":
							rebateTrans.Rebate = detailsStr
						case "ccy":
							rebateTrans.RebateCcy = detailsStr
						}
						return err
					})
					if err != nil {
						return
					}
					details.RebateTrans = append(details.RebateTrans, *rebateTrans)
				})
			case "triggerParams":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						switch string(key) {
						case "triggerAction":
							triggerParams.TriggerAction = detailsStr
						case "triggerStrategy":
							triggerParams.TriggerStrategy = detailsStr
						case "delaySeconds":
							triggerParams.DelaySeconds = detailsStr
						case "triggerTime":
							triggerParams.TriggerTime = detailsStr
						case "triggerType":
							triggerParams.TriggerType = detailsStr
						case "timeframe":
							triggerParams.Timeframe = detailsStr
						case "thold":
							triggerParams.Thold = detailsStr
						case "triggerCond":
							triggerParams.TriggerCond = detailsStr
						case "timePeriod":
							triggerParams.TimePeriod = detailsStr
						case "triggerPx":
							triggerParams.TriggerPx = detailsStr
						case "stopType":
							triggerParams.StopType = detailsStr
						}
						return err
					})
					if err != nil {
						return
					}
					details.TriggerParams = append(details.TriggerParams, *triggerParams)
				})
			case "algoId":
				details.AlgoId = detailsStr
			case "algoClOrdId":
				details.AlgoClOrdId = detailsStr
			case "instType":
				details.InstType = detailsStr
			case "instId":
				details.InstId = detailsStr
			case "cTime":
				details.CTime = detailsStr
			case "uTime":
				details.UTime = detailsStr
			case "algoOrdType":
				details.AlgoOrdType = detailsStr
			case "state":
				details.State = detailsStr
			case "maxPx":
				details.MaxPx = detailsStr
			case "minPx":
				details.MinPx = detailsStr
			case "gridNum":
				details.GridNum = detailsStr
			case "runType":
				details.RunType = detailsStr
			case "tpTriggerPx":
				details.TpTriggerPx = detailsStr
			case "slTriggerPx":
				details.SlTriggerPx = detailsStr
			case "tradeNum":
				details.TradeNum = detailsStr
			case "arbitrageNum":
				details.ArbitrageNum = detailsStr
			case "singleAmt":
				details.SingleAmt = detailsStr
			case "perMinProfitRate":
				details.PerMinProfitRate = detailsStr
			case "perMaxProfitRate":
				details.PerMaxProfitRate = detailsStr
			case "runPx":
				details.RunPx = detailsStr
			case "totalPnl":
				details.TotalPnl = detailsStr
			case "pnlRatio":
				details.PnlRatio = detailsStr
			case "investment":
				details.Investment = detailsStr
			case "gridProfit":
				details.GridProfit = detailsStr
			case "floatProfit":
				details.FloatProfit = detailsStr
			case "totalAnnualizedRate":
				details.TotalAnnualizedRate = detailsStr
			case "annualizedRate":
				details.AnnualizedRate = detailsStr
			case "cancelType":
				details.CancelType = detailsStr
			case "stopType":
				details.StopType = detailsStr
			case "activeOrdNum":
				details.ActiveOrdNum = detailsStr
			case "quoteSz":
				details.QuoteSz = detailsStr
			case "baseSz":
				details.BaseSz = detailsStr
			case "curQuoteSz":
				details.CurQuoteSz = detailsStr
			case "curBaseSz":
				details.CurBaseSz = detailsStr
			case "profit":
				details.Profit = detailsStr
			case "stopResult":
				details.StopResult = detailsStr
			case "direction":
				details.Direction = detailsStr
			case "basePos":
				details.BasePos = detailsStr
			case "sz":
				details.Sz = detailsStr
			case "lever":
				details.Lever = detailsStr
			case "actualLever":
				details.ActualLever = detailsStr
			case "liqPx":
				details.LiqPx = detailsStr
			case "uly":
				details.Uly = detailsStr
			case "instFamily":
				details.InstFamily = detailsStr
			case "ordFrozen":
				details.OrdFrozen = detailsStr
			case "availEq":
				details.AvailEq = detailsStr
			case "eq":
				details.Eq = detailsStr
			case "tag":
				details.Tag = detailsStr
			case "profitSharingRatio":
				details.ProfitSharingRatio = detailsStr
			case "copyType":
				details.CopyType = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return *details, err
}

func (un *RespUnmarshaler) UnmarshalPostPlaceGridAlgoOrder(data []byte) (model.PlaceGridAlgoOrderResponse, error) {
	var details = new(model.PlaceGridAlgoOrderResponse)

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
			case "tag":
				details.Tag = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return *details, err
}

func (un *RespUnmarshaler) UnmarshalPlaceOrder(respPlaceOrderData []byte) (model.PlaceOrderResponse, error) {
	var placeOrderResponse = new(model.PlaceOrderResponse)
	var placeOrderResponseData = new(model.PlaceOrderResponseData)

	_, err := jsonparser.ArrayEach(respPlaceOrderData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "code":
				placeOrderResponse.Code = detailsStr
			case "msg":
				placeOrderResponse.Msg = detailsStr
			case "inTime":
				placeOrderResponse.InTime = detailsStr
			case "outTime":
				placeOrderResponse.OutTime = detailsStr
			case "data":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						switch string(key) {
						case "ordId":
							placeOrderResponseData.OrdId = detailsStr
						case "clOrdId":
							placeOrderResponseData.ClOrdId = detailsStr
						case "tag":
							placeOrderResponseData.Tag = detailsStr
						case "sCode":
							placeOrderResponseData.SCode = detailsStr
						case "sMsg":
							logger.Info("++++++++++ sMsg - ", detailsStr)
							placeOrderResponseData.SMsg = detailsStr
						}
						return err
					})
					if err != nil {
						return
					}
					placeOrderResponse.Data = append(placeOrderResponse.Data, *placeOrderResponseData)
				})
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return *placeOrderResponse, err
}

func (un *RespUnmarshaler) UnmarshalAmendOrderResponse(respAmendOrderData []byte) (model.AmendOrderResponse, error) {
	var amendOrderResponse = new(model.AmendOrderResponse)
	var amendOrderResponseData = new(model.AmendOrderResponseData)

	_, err := jsonparser.ArrayEach(respAmendOrderData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "code":
				amendOrderResponse.Code = detailsStr
			case "msg":
				amendOrderResponse.Msg = detailsStr
			case "inTime":
				amendOrderResponse.InTime = detailsStr
			case "outTime":
				amendOrderResponse.OutTime = detailsStr
			case "data":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						switch string(key) {
						case "ordId":
							amendOrderResponseData.OrdId = detailsStr
						case "clOrdId":
							amendOrderResponseData.ClOrdId = detailsStr
						case "tag":
							amendOrderResponseData.ReqId = detailsStr
						case "sCode":
							amendOrderResponseData.SCode = detailsStr
						case "sMsg":
							logger.Info("++++++++++ sMsg - ", detailsStr)
							amendOrderResponseData.SMsg = detailsStr
						}
						return err
					})
					if err != nil {
						return
					}
					amendOrderResponse.Data = append(amendOrderResponse.Data, *amendOrderResponseData)
				})
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return *amendOrderResponse, err
}

func (un *RespUnmarshaler) UnmarshalPostStopGridAlgoOrder(data []byte) (model.StopGridAlgoOrderResponse, error) {
	var details = new(model.StopGridAlgoOrderResponse)

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
			case "tag":
				details.Tag = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return *details, err
}
