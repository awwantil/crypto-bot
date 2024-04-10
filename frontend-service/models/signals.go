package models

import (
	"github.com/aidarkhanov/nanoid"
	"gorm.io/gorm"
	u "okx-bot/frontend-service/utils"
)

type TradingViewSignal struct {
	gorm.Model
	IdOrder                string `json:"idOrder"`
	Action                 string `json:"action"`
	MarketPosition         string `json:"marketPosition" sql:"-"`
	PrevMarketPosition     string `json:"prevMarketPosition"`
	MarketPositionSize     string `json:"marketPositionSize"`
	PrevMarketPositionSize string `json:"prevMarketPositionSize"`
	Instrument             string `json:"instrument"`
	SignalToken            string `json:"signalToken"`
	Timestamp              string `json:"timestamp"`
	Amount                 string `json:"amount"`
}

type SignalObject struct {
	gorm.Model
	SignalId     string `gorm:"uniqueIndex"`
	NameToken    string `json:"nameToken"`
	TimeInterval string `json:"timeInterval"`
}

func (tradingViewSignal *TradingViewSignal) Save() map[string]interface{} {

	GetDB().Create(tradingViewSignal)

	response := u.Message(true, "TradingViewSignal has been saved")
	response["tradingViewSignal"] = tradingViewSignal
	return response
}

func (signalObject *SignalObject) Create(nameToken string, interval string) map[string]interface{} {
	signalObject.SignalId = nanoid.New()
	signalObject.NameToken = nameToken
	signalObject.TimeInterval = interval

	GetDB().Create(signalObject)

	response := u.Message(true, "SignalObject has been created")
	response["signalObject"] = signalObject
	return response
}
