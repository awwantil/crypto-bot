package model

// https://www.okx.com/docs-v5/en/#order-book-trading-signal-bot-trading-post-create-signal
// POST /api/v5/tradingBot/signal/create-signal
type CreateSignalRequest struct {
	SignalChanName string `json:"signalChanName,omitempty"`
	SignalChanDesc string `json:"signalChanDesc,omitempty"`
}

type CreateSignalResponse struct {
	SignalChanId    string `json:"signalChanId,omitempty"`
	SignalChanToken string `json:"signalChanToken,omitempty"`
}

// https://www.okx.com/docs-v5/en/#order-book-trading-signal-bot-trading-post-create-signal-bot
// POST /api/v5/tradingBot/signal/order-algo
type CreateSignalBotRequest struct {
	SignalChanId      string                `json:"signalChanId,omitempty"`
	Lever             string                `json:"lever,omitempty"`
	InvestAmt         string                `json:"investAmt,omitempty"`
	SubOrdType        string                `json:"subOrdType,omitempty"`
	InstIds           []string              `json:"instIds,omitempty"`
	IncludeAll        string                `json:"includeAll,omitempty"`
	EntrySettingParam EntrySettingParamData `json:"entrySettingParam,omitempty"`
	ExitSettingParam  ExitSettingParamData  `json:"exitSettingParam,omitempty"`
}

type EntrySettingParamData struct {
	AllowMultipleEntry string `json:"allowMultipleEntry,omitempty"`
	EntryType          string `json:"entryType,omitempty"`
	Amt                string `json:"amt,omitempty"`
	Ratio              string `json:"ratio,omitempty"`
}

type ExitSettingParamData struct {
	TpSlType string `json:"tpSlType,omitempty"`
	TpPct    string `json:"tpPct,omitempty"`
	SlPct    string `json:"slPct,omitempty"`
}

type CreateSignalBotResponse struct {
	AlgoId      string `json:"algoId,omitempty"`
	AlgoClOrdId string `json:"algoClOrdId,omitempty"`
	SCode       string `json:"sCode,omitempty"`
	SMsg        string `json:"sMsg,omitempty"`
}

// https://www.okx.com/docs-v5/en/#order-book-trading-signal-bot-trading-post-cancel-signal-bots
// POST /api/v5/tradingBot/signal/stop-order-algo
type CancelSignalBotRequest struct {
	AlgoId string `json:"algoId,omitempty"`
}

type CancelSignalBotResponse struct {
	AlgoId string `json:"algoId,omitempty"`
	SCode  string `json:"sCode,omitempty"`
	SMsg   string `json:"sMsg,omitempty"`
}

// https://www.okx.com/docs-v5/en/#order-book-trading-signal-bot-trading-post-place-sub-order
// POST /api/v5/tradingBot/signal/sub-order
type PlaceSubOrderSignalBotRequest struct {
	InstId  string `json:"instId,omitempty"`
	AlgoId  string `json:"algoId,omitempty"`
	Side    string `json:"side,omitempty"`
	OrdType string `json:"ordType,omitempty"`
	Sz      string `json:"sz,omitempty"`
	Px      string `json:"px,omitempty"`
}

type PlaceSubOrderSignalBotResponse struct {
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data string `json:"data,omitempty"`
}

// https://www.okx.com/docs-v5/en/#order-book-trading-signal-bot-trading-post-cancel-sub-order
// POST /api/v5/tradingBot/signal/cancel-sub-order
type CancelSubOrderSignalBotRequest struct {
	InstId      string `json:"instId,omitempty"`
	AlgoId      string `json:"algoId,omitempty"`
	SignalOrdId string `json:"signalOrdId,omitempty"`
}

type CancelSubOrderSignalBotResponse struct {
	Code string               `json:"code,omitempty"`
	Msg  string               `json:"msg,omitempty"`
	Data []CancelSubOrderData `json:"data,omitempty"`
}

type CancelSubOrderData struct {
	SignalOrdId string `json:"signalOrdId,omitempty"`
	SCode       string `json:"sCode,omitempty"`
	SMsg        string `json:"sMsg,omitempty"`
}

// https://www.okx.com/docs-v5/en/#order-book-trading-signal-bot-trading-post-close-position
// POST /api/v5/tradingBot/signal/close-position
type ClosePositionSignalBotRequest struct {
	InstId string `json:"instId,omitempty"`
	AlgoId string `json:"algoId,omitempty"`
}

type ClosePositionSignalBotResponse struct {
	AlgoId string `json:"algoId,omitempty"`
}

// https://www.okx.com/docs-v5/en/#order-book-trading-signal-bot-trading-get-signal-bot-sub-orders
// GET /api/v5/tradingBot/signal/sub-orders
type GetSubOrdersSignalBotRequest struct {
	AlgoId      string `json:"algoId,omitempty"`
	AlgoOrdType string `json:"algoOrdType,omitempty"`
	SignalOrdId string `json:"signalOrdId,omitempty"`
	State       string `json:"state,omitempty"`
}

type GetSubOrdersSignalBotResponse struct {
	Adl         string `json:"adl,omitempty"`
	AlgoClOrdId string `json:"algoClOrdId,omitempty"`
	AlgoId      string `json:"algoId,omitempty"`
	AvgPx       string `json:"avgPx,omitempty"`
	CTime       string `json:"cTime,omitempty"`
	Ccy         string `json:"ccy,omitempty"`
	Imr         string `json:"imr,omitempty"`
	InstId      string `json:"instId,omitempty"`
	InstType    string `json:"instType,omitempty"`
	Last        string `json:"last,omitempty"`
	Lever       string `json:"lever,omitempty"`
	LiqPx       string `json:"liqPx,omitempty"`
	MarkPx      string `json:"markPx,omitempty"`
	MgnMode     string `json:"mgnMode,omitempty"`
	MgnRatio    string `json:"mgnRatio,omitempty"`
	Mmr         string `json:"mmr,omitempty"`
	NotionalUsd string `json:"notionalUsd,omitempty"`
	OrdId       string `json:"ordId,omitempty"`
	Px          string `json:"px,omitempty"`
	Pos         string `json:"pos,omitempty"`
	PosSide     string `json:"posSide,omitempty"`
	UTime       string `json:"uTime,omitempty"`
	Upl         string `json:"upl,omitempty"`
	UplRatio    string `json:"uplRatio,omitempty"`
}

// https://www.okx.com/docs-v5/en/#order-book-trading-signal-bot-trading-get-signals
// GET /api/v5/tradingBot/signal/signals
type GetSignalsRequest struct {
	SignalSourceType string `json:"signalSourceType,omitempty"`
}

type GetSignalsResponse struct {
	Code string                   `json:"code,omitempty"`
	Msg  string                   `json:"msg,omitempty"`
	Data []GetSignalsResponseData `json:"data,omitempty"`
}

type GetSignalsResponseData struct {
	SignalChanId     string `json:"signalChanId,omitempty"`
	SignalChanName   string `json:"signalChanName,omitempty"`
	SignalChanDesc   string `json:"signalChanDesc,omitempty"`
	SignalChanToken  string `json:"signalChanToken,omitempty"`
	SignalSourceType string `json:"signalSourceType,omitempty"`
}

// https://www.okx.com/docs-v5/en/#order-book-trading-signal-bot-trading-get-active-signal-bot
// GET /api/v5/tradingBot/signal/orders-algo-pending
type GetActiveSignalBotRequest struct {
	AlgoOrdType string `json:"algoOrdType,omitempty"`
	AlgoId      string `json:"algoId,omitempty"`
	After       string `json:"after,omitempty"`
	Before      string `json:"before,omitempty"`
}

type GetActiveSignalBotResponse struct {
	AvailBal     string `json:"availBal,omitempty"`
	FrozenBal    string `json:"frozenBal,omitempty"`
	SignalChanId string `json:"signalChanId,omitempty"`
	InvestAmt    string `json:"investAmt,omitempty"`
	Lever        string `json:"lever,omitempty"`
	FloatPnl     string `json:"floatPnl,omitempty"`
	RealizedPnl  string `json:"realizedPnl,omitempty"`
	TotalPnl     string `json:"totalPnl,omitempty"`
	TotalEq      string `json:"totalEq,omitempty"`
	CTime        string `json:"cTime,omitempty"`
	UTime        string `json:"uTime,omitempty"`
}
