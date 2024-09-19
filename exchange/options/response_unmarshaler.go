package options

import "okx-bot/exchange/model"

type ResponseUnmarshaler func([]byte, interface{}) error
type GetTickerResponseUnmarshaler func([]byte) (*model.Ticker, error)
type GetDepthResponseUnmarshaler func([]byte) (*model.Depth, error)
type GetKlineResponseUnmarshaler func([]byte) ([]model.Kline, error)
type CreateOrderResponseUnmarshaler func([]byte) (*model.Order, error)
type SetLeverageResponseUnmarshaler func([]byte) (*model.SetLeverageResponse, error)
type GetOrderInfoResponseUnmarshaler func([]byte) (*model.OrderInfoResponse, error)
type GetPendingOrdersResponseUnmarshaler func([]byte) ([]model.OrderInfoResponse, error)
type CancelOrderResponseUnmarshaler func([]byte) (model.CancelOrderResponse, error)
type ClosePositionsResponseUnmarshaler func([]byte) (model.ClosePositionResponse, error)
type GetHistoryOrdersResponseUnmarshaler func([]byte) ([]model.OrderInfoResponse, error)
type GetAccountResponseUnmarshaler func([]byte) (map[string]model.Account, error)
type GetPositionsResponseUnmarshaler func([]byte) ([]model.FuturesPosition, error)
type GetPositionsHistoryResponseUnmarshaler func([]byte) ([]model.FuturesPositionHistory, error)
type GetAccountBalanceResponseUnmarshaler func([]byte) (model.BalanceResponse, error)
type GetFuturesAccountResponseUnmarshaler func([]byte) (map[string]model.FuturesAccount, error)
type GetExchangeInfoResponseUnmarshaler func([]byte) (map[string]model.CurrencyPair, error)
type GetCompMinInvestResponseUnmarshaler func([]byte) (model.ComputeMinInvestmentResponse, error)
type GetAlgoOrderDetailsResponseUnmarshaler func([]byte) (model.GridAlgoOrderDetailsResponse, error)
type PlaceGridAlgoOrderResponseUnmarshaler func([]byte) (model.PlaceGridAlgoOrderResponse, error)
type StopGridAlgoOrderResponseUnmarshaler func([]byte) (model.StopGridAlgoOrderResponse, error)
type PlaceOrderResponseUnmarshaler func([]byte) (model.PlaceOrderResponse, error)
type AmendOrderResponseUnmarshaler func([]byte) (model.AmendOrderResponse, error)

type PostCreateSignalUnmarshaler func([]byte) (*model.CreateSignalResponse, error)
type PostCreateSignalBotUnmarshaler func([]byte) (*model.CreateSignalBotResponse, error)
type PostCancelSignalBotUnmarshaler func([]byte) (*model.CancelSignalBotResponse, error)
type PostPlaceSubOrderSignalBotUnmarshaler func([]byte) (*model.PlaceSubOrderSignalBotResponse, error)
type PostCancelSubOrderSignalBotUnmarshaler func([]byte) (*model.CancelSubOrderSignalBotResponse, error)
type PostClosePositionSignalBotUnmarshaler func([]byte) (*model.ClosePositionSignalBotResponse, error)
type GetSubOrdersSignalBotUnmarshaler func([]byte) (*model.GetSubOrdersSignalBotResponse, error)

type UnmarshalerOptions struct {
	ResponseUnmarshaler                    ResponseUnmarshaler
	TickerUnmarshaler                      GetTickerResponseUnmarshaler
	DepthUnmarshaler                       GetDepthResponseUnmarshaler
	KlineUnmarshaler                       GetKlineResponseUnmarshaler
	CreateOrderResponseUnmarshaler         CreateOrderResponseUnmarshaler
	SetLeverageResponseUnmarshaler         SetLeverageResponseUnmarshaler
	GetOrderInfoResponseUnmarshaler        GetOrderInfoResponseUnmarshaler
	GetPendingOrdersResponseUnmarshaler    GetPendingOrdersResponseUnmarshaler
	GetHistoryOrdersResponseUnmarshaler    GetHistoryOrdersResponseUnmarshaler
	CancelOrderResponseUnmarshaler         CancelOrderResponseUnmarshaler
	ClosePositionsResponseUnmarshaler      ClosePositionsResponseUnmarshaler
	AmendOrderResponseUnmarshaler          AmendOrderResponseUnmarshaler
	GetAccountResponseUnmarshaler          GetAccountResponseUnmarshaler
	GetPositionsResponseUnmarshaler        GetPositionsResponseUnmarshaler
	GetPositionsHistoryResponseUnmarshaler GetPositionsHistoryResponseUnmarshaler
	GetAccountBalanceResponseUnmarshaler   GetAccountBalanceResponseUnmarshaler
	GetFuturesAccountResponseUnmarshaler   GetFuturesAccountResponseUnmarshaler
	GetExchangeInfoResponseUnmarshaler     GetExchangeInfoResponseUnmarshaler
	GetCompMinInvestResponseUnmarshaler    GetCompMinInvestResponseUnmarshaler
	GetAlgoOrderDetailsResponseUnmarshaler GetAlgoOrderDetailsResponseUnmarshaler
	PlaceGridAlgoOrderResponseUnmarshaler  PlaceGridAlgoOrderResponseUnmarshaler
	StopGridAlgoOrderResponseUnmarshaler   StopGridAlgoOrderResponseUnmarshaler
	PlaceOrderResponseUnmarshaler          PlaceOrderResponseUnmarshaler

	PostCreateSignalUnmarshaler            PostCreateSignalUnmarshaler
	PostCreateSignalBotUnmarshaler         PostCreateSignalBotUnmarshaler
	PostCancelSignalBotUnmarshaler         PostCancelSignalBotUnmarshaler
	PostPlaceSubOrderSignalBotUnmarshaler  PostPlaceSubOrderSignalBotUnmarshaler
	PostCancelSubOrderSignalBotUnmarshaler PostCancelSubOrderSignalBotUnmarshaler
	PostClosePositionSignalBotUnmarshaler  PostClosePositionSignalBotUnmarshaler
	GetSubOrdersSignalBotUnmarshaler       GetSubOrdersSignalBotUnmarshaler
}

type UnmarshalerOption func(options *UnmarshalerOptions)

func WithResponseUnmarshaler(unmarshaler ResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.ResponseUnmarshaler = unmarshaler
	}
}

func WithTickerUnmarshaler(unmarshaler GetTickerResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.TickerUnmarshaler = unmarshaler
	}
}

func WithDepthUnmarshaler(unmarshaler GetDepthResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.DepthUnmarshaler = unmarshaler
	}
}

func WithKlineUnmarshaler(unmarshaler GetKlineResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.KlineUnmarshaler = unmarshaler
	}
}

func WithGetOrderInfoResponseUnmarshaler(unmarshaler GetOrderInfoResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetOrderInfoResponseUnmarshaler = unmarshaler
	}
}

func WithCreateOrderResponseUnmarshaler(unmarshaler CreateOrderResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.CreateOrderResponseUnmarshaler = unmarshaler
	}
}

func WithGetPendingOrdersResponseUnmarshaler(unmarshaler GetPendingOrdersResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetPendingOrdersResponseUnmarshaler = unmarshaler
	}
}

func WithCancelOrderResponseUnmarshaler(unmarshaler CancelOrderResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.CancelOrderResponseUnmarshaler = unmarshaler
	}
}

func WithGetHistoryOrdersResponseUnmarshaler(unmarshaler GetHistoryOrdersResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetHistoryOrdersResponseUnmarshaler = unmarshaler
	}
}

func WithGetAccountResponseUnmarshaler(unmarshaler GetAccountResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetAccountResponseUnmarshaler = unmarshaler
	}
}

func WithGetPositionsResponseUnmarshaler(unmarshaler GetPositionsResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetPositionsResponseUnmarshaler = unmarshaler
	}
}

func WithGetPositionsHistoryResponseUnmarshaler(unmarshaler GetPositionsHistoryResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetPositionsHistoryResponseUnmarshaler = unmarshaler
	}
}

func WithGetFuturesAccountResponseUnmarshaler(unmarshaler GetFuturesAccountResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetFuturesAccountResponseUnmarshaler = unmarshaler
	}
}

func WithGetExchangeInfoResponseUnmarshaler(unmarshaler GetExchangeInfoResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetExchangeInfoResponseUnmarshaler = unmarshaler
	}
}
