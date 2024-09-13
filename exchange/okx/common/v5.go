package common

import (
	"encoding/json"
	"okx-bot/exchange/options"
)

type OKxV5 struct {
	UriOpts       options.UriOptions
	UnmarshalOpts options.UnmarshalerOptions
}

type BaseResp struct {
	Code int             `json:"code,string"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func New() *OKxV5 {
	unmarshaler := new(RespUnmarshaler)

	f := &OKxV5{
		UriOpts: options.UriOptions{
			Endpoint:               "https://www.okx.com",
			KlineUri:               "/api/v5/market/candles",
			TickerUri:              "/api/v5/market/ticker",
			DepthUri:               "/api/v5/market/books",
			NewOrderUri:            "/api/v5/trade/order",
			SetLeverageUri:         "/api/v5/account/set-leverage",
			GetOrderUri:            "/api/v5/trade/order",
			AmendOrderUri:          "/api/v5/trade/amend-order",
			GetHistoryOrdersUri:    "/api/v5/trade/orders-history",
			GetPendingOrdersUri:    "/api/v5/trade/orders-pending",
			CancelOrderUri:         "/api/v5/trade/cancel-order",
			ClosePositionsUri:      "/api/v5/trade/close-position",
			GetAccountUri:          "/api/v5/account/balance",
			GetPositionsUri:        "/api/v5/account/positions",
			GetPositionsHistoryUri: "/api/v5/account/positions-history",
			GetAccountBalanceUri:   "/api/v5/account/balance",
			GetExchangeInfoUri:     "/api/v5/public/instruments",

			PostPlaceGridAlgoOrderUri: "/api/v5/tradingBot/grid/order-algo",
			PostStopGridAlgoOrderUri:  "/api/v5/tradingBot/grid/stop-order-algo",
			PostComputeMinInvestment:  "/api/v5/tradingBot/grid/min-investment",
			GetAlgoOrderDetails:       "/api/v5/tradingBot/grid/orders-algo-details",

			PostCreateSignalUri:            "/api/v5/tradingBot/signal/create-signal",
			PostCreateSignalBotUri:         "/api/v5/tradingBot/signal/order-algo",
			PostCancelSignalBotUri:         "/api/v5/tradingBot/signal/stop-order-algo",
			PostPlaceSubOrderSignalBotUri:  "/api/v5/tradingBot/signal/sub-order",
			PostCancelSubOrderSignalBotUri: "/api/v5/tradingBot/signal/cancel-sub-order",
		},
		UnmarshalOpts: options.UnmarshalerOptions{
			ResponseUnmarshaler:                    unmarshaler.UnmarshalResponse,
			KlineUnmarshaler:                       unmarshaler.UnmarshalGetKlineResponse,
			TickerUnmarshaler:                      unmarshaler.UnmarshalTicker,
			DepthUnmarshaler:                       unmarshaler.UnmarshalDepth,
			CreateOrderResponseUnmarshaler:         unmarshaler.UnmarshalCreateOrderResponse,
			SetLeverageResponseUnmarshaler:         unmarshaler.UnmarshalSetLeverageResponse,
			GetPendingOrdersResponseUnmarshaler:    unmarshaler.UnmarshalGetPendingOrdersResponse,
			GetHistoryOrdersResponseUnmarshaler:    unmarshaler.UnmarshalGetHistoryOrdersResponse,
			CancelOrderResponseUnmarshaler:         unmarshaler.UnmarshalCancelOrderResponse,
			ClosePositionsResponseUnmarshaler:      unmarshaler.UnmarshalClosePositionsResponse,
			GetOrderInfoResponseUnmarshaler:        unmarshaler.UnmarshalGetOrderInfoResponse,
			AmendOrderResponseUnmarshaler:          unmarshaler.UnmarshalAmendOrderResponse,
			GetAccountResponseUnmarshaler:          unmarshaler.UnmarshalGetAccountResponse,
			GetPositionsResponseUnmarshaler:        unmarshaler.UnmarshalGetPositionsResponse,
			GetFuturesAccountResponseUnmarshaler:   unmarshaler.UnmarshalGetFuturesAccountResponse,
			GetPositionsHistoryResponseUnmarshaler: unmarshaler.UnmarshalGetPositionsHisotoryResponse,
			GetAccountBalanceResponseUnmarshaler:   unmarshaler.UnmarshalGetAccountBalanceResponse,
			GetExchangeInfoResponseUnmarshaler:     unmarshaler.UnmarshalGetExchangeInfoResponse,
			GetCompMinInvestResponseUnmarshaler:    unmarshaler.UnmarshalGetComputeMinInvestmentResponse,
			GetAlgoOrderDetailsResponseUnmarshaler: unmarshaler.UnmarshalGetAlgoOrderDetailsResponse,
			PlaceGridAlgoOrderResponseUnmarshaler:  unmarshaler.UnmarshalPostPlaceGridAlgoOrder,
			StopGridAlgoOrderResponseUnmarshaler:   unmarshaler.UnmarshalPostStopGridAlgoOrder,
			PlaceOrderResponseUnmarshaler:          unmarshaler.UnmarshalPlaceOrder,

			PostCreateSignalUnmarshaler:            unmarshaler.UnmarshalCreateSignal,
			PostCreateSignalBotUnmarshaler:         unmarshaler.UnmarshalCreateSignalBot,
			PostCancelSignalBotUnmarshaler:         unmarshaler.UnmarshalCancelSignalBot,
			PostPlaceSubOrderSignalBotUnmarshaler:  unmarshaler.UnmarshalPlaceSubOrderSignalBot,
			PostCancelSubOrderSignalBotUnmarshaler: unmarshaler.UnmarshalCancelSubOrderSignalBot,
		},
	}

	return f
}

func (okx *OKxV5) WithUriOption(opts ...options.UriOption) *OKxV5 {
	for _, opt := range opts {
		opt(&okx.UriOpts)
	}
	return okx
}

func (okx *OKxV5) WithUnmarshalOption(opts ...options.UnmarshalerOption) *OKxV5 {
	for _, opt := range opts {
		opt(&okx.UnmarshalOpts)
	}
	return okx
}

func (okx *OKxV5) NewPrvApi(opts ...options.ApiOption) *Prv {
	api := NewPrvApi(opts...)
	api.OKxV5 = okx
	return api
}
