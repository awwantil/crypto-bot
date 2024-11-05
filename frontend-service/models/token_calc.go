package models

const (
	BASE_CURRENCY = "USDT"
)

const (
	ADA  = "ADA"
	ATOM = "ATOM"
	AVAX = "AVAX"
	APT  = "APT"

	BTC = "BTC"
	BCH = "BCH"
	BNB = "BNB"

	CRO = "CRO"

	ETH = "ETH"
	ETC = "ETC"

	DOGE = "DOGE"
	DOT  = "DOT"

	FIL = "FIL"

	ICP = "ICP"
	IMX = "IMX"

	LTC  = "LTC"
	NEAR = "NEAR"

	SOL  = "SOL"
	SHIB = "SHIB"

	TON = "TON"

	USD  = "USD"
	USDT = "USDT"
	UNI  = "UNI"

	VET = "VET"

	XRP = "XRP"
	XLM = "XLM"
)

type CalcPriceData struct {
	DemoStep      float64
	DemoMinAmount float64
	DemoPrecision int
	ProdStep      float64
	ProdMinAmount float64
	ProdPrecision int
}

// for calcPx
// https://www.okx.com/ru/trade-market/info/swap
var PriceData = map[string]CalcPriceData{
	ETH:  {DemoStep: 0.1, DemoMinAmount: 0.001, DemoPrecision: 1, ProdStep: 0.1, ProdMinAmount: 0.01, ProdPrecision: 1},
	XRP:  {DemoStep: 0.1, DemoMinAmount: 10, DemoPrecision: 1, ProdStep: 0.1, ProdMinAmount: 10, ProdPrecision: 1},
	SOL:  {DemoStep: 0.01, DemoMinAmount: 0.001, DemoPrecision: 2, ProdStep: 0.01, ProdMinAmount: 0.01, ProdPrecision: 2},
	ADA:  {DemoStep: 0.1, DemoMinAmount: 0.1, DemoPrecision: 1, ProdStep: 0.1, ProdMinAmount: 10, ProdPrecision: 1},
	DOGE: {DemoStep: 0.1, DemoMinAmount: 1, DemoPrecision: 1, ProdStep: 0.1, ProdMinAmount: 100, ProdPrecision: 1},
	LTC:  {DemoStep: 0.1, DemoMinAmount: 0.1, DemoPrecision: 1, ProdStep: 0.1, ProdMinAmount: 0.1, ProdPrecision: 1},
}
