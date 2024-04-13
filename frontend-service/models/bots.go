package models

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	u "okx-bot/frontend-service/utils"
	"time"
)

type BotStatus int

const (
	Created BotStatus = iota + 1
	Waiting
	MakingDeal
	Stopped
	Error
	Pause
)

type Bot struct {
	gorm.Model
	ID            uint      `gorm:"primary_key"`
	StartTime     time.Time `json:"startTime"`
	Status        BotStatus `json:"status"`
	InitialAmount float64   `json:"initialAmount"`
	CurrentAmount float64   `json:"currentAmount"`
	SignalRefer   uint
	//Deals     []Deal       `json:"deals" gorm:"foreignKey:UserRefer"`
}

type BotRequest struct {
	gorm.Model
	InitialAmount float64 `json:"initialAmount"`
	CodeSignalId  string  `json:"codeSignalId"`
}

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.models.base",
	})
)

func (index BotStatus) ToString() string {
	return [...]string{"Created", "Waiting", "MakingDeal", "Stopped", "Error", "Pause"}[index-1]
}

func (index BotStatus) EnumIndex() int {
	return int(index)
}

func (bot *Bot) Create(codeId string, initialAmount float64) map[string]interface{} {
	var signal = Signal{Code: codeId}
	db.First(&signal)

	logger.Infoln("signal", signal)

	bot.Status = Created
	bot.InitialAmount = initialAmount
	//GetDB().Create(bot)

	signal.Bot = append(signal.Bot, *bot)
	db.Save(&signal)

	response := u.Message(true, "Bot has been created")
	response["bot"] = bot
	return response
}
