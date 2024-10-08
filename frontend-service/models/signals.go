package models

import (
	"fmt"
	"github.com/aidarkhanov/nanoid"
	"gorm.io/gorm"
	u "okx-bot/frontend-service/utils"
)

type TradingViewSignalReceive struct {
	gorm.Model
	IdSignal               uint   `gorm:"primary_key"`
	Id                     string `json:"id"`
	Action                 string `json:"action"`
	MarketPosition         string `json:"marketPosition"`
	PrevMarketPosition     string `json:"prevMarketPosition" sql:"-"`
	MarketPositionSize     string `json:"marketPositionSize" sql:"-"`
	PrevMarketPositionSize string `json:"prevMarketPositionSize" sql:"-"`
	Instrument             string `json:"instrument"`
	SignalToken            string `json:"signalToken"`
	Timestamp              string `json:"timestamp"`
	InvestmentType         string `json:"investmentType"`
	Amount                 string `json:"amount"`
}

type Signal struct {
	gorm.Model
	ID           uint   `gorm:"primary_key"`
	StrategyName string `gorm:"strategyName"`
	StrategyDesc string `gorm:"strategyDesc"`
	Code         string `gorm:"uniqueIndex"`
	NameToken    string `json:"nameToken"`
	TimeInterval string `json:"timeInterval"`
	Bots         []Bot  `gorm:"foreignKey:SignalRefer"`
}

func (tradingViewSignal *TradingViewSignalReceive) Save() map[string]interface{} {

	GetDB().Create(tradingViewSignal)

	response := u.Message(true, "TradingViewSignal has been saved")
	response["tradingViewSignal"] = tradingViewSignal
	return response
}

func (signal *Signal) Create() map[string]interface{} {
	signal.Code = nanoid.New()
	GetDB().Create(signal)

	response := u.Message(true, "Signal has been created")
	response["signal"] = signal
	return response
}

func GetAllSignals() *[]Signal {
	findSignals := &[]Signal{}

	err := GetDB().Table("signals").Find(&findSignals).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return findSignals
}

func GetBots(signalCodeId string) []Bot {
	findSignal := &Signal{}

	err := GetDB().Table("signals").Where("code = ?", signalCodeId).Preload("Bots").Find(&findSignal).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return findSignal.Bots
}

func FindSignalByCode(signalCodeId string) (*Signal, error) {
	findSignal := &Signal{}
	err := GetDB().Table("signals").Where("code = ?", signalCodeId).Find(&findSignal).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return findSignal, nil
}

func FindSignalById(signalId uint) (*Signal, error) {
	findSignal := &Signal{}
	err := GetDB().Table("signals").Where("id = ?", signalId).Find(&findSignal).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return findSignal, nil
}
