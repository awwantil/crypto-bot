package models

import (
	"github.com/aidarkhanov/nanoid"
	"gorm.io/gorm"
	u "okx-bot/frontend-service/utils"
)

type TradingViewSignalReceive struct {
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

type Signal struct {
	gorm.Model
	ID           uint   `gorm:"primary_key"`
	Code         string `gorm:"uniqueIndex"`
	NameToken    string `json:"nameToken"`
	TimeInterval string `json:"timeInterval"`
	Bot          []Bot  `gorm:"foreignKey:SignalRefer"`
}

func (tradingViewSignal *TradingViewSignalReceive) Save() map[string]interface{} {

	GetDB().Create(tradingViewSignal)

	response := u.Message(true, "TradingViewSignal has been saved")
	response["tradingViewSignal"] = tradingViewSignal
	return response
}

func (signal *Signal) Create(nameToken string, interval string) map[string]interface{} {
	signal.Code = nanoid.New()
	signal.NameToken = nameToken
	signal.TimeInterval = interval

	GetDB().Create(signal)

	response := u.Message(true, "Signal has been created")
	response["signal"] = signal
	return response
}
