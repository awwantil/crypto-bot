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
	Deals         []Deal `json:"deals" gorm:"foreignKey:BotRefer"`
}

type BotCreateRequest struct {
	InitialAmount float64 `json:"initialAmount"`
	CodeSignalId  string  `json:"codeSignalId"`
}

type BotWithIdRequest struct {
	Id string `json:"id"`
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
	db.Where("code = ?", codeId).First(&signal)

	logger.Infoln("signal", signal)

	bot.Status = Created
	bot.InitialAmount = initialAmount
	bot.StartTime = time.Now()

	signal.Bots = append(signal.Bots, *bot)
	db.Save(&signal)

	response := u.Message(true, "Bot has been created")
	response["bot"] = bot
	return response
}

func (bot *Bot) Delete() map[string]interface{} {

	db.Delete(&bot)

	response := u.Message(true, "Bot has been deleted")
	return response
}

func Find(botIdRequest BotWithIdRequest) Bot {
	var foundBot = Bot{}
	db.Where("id = ?", botIdRequest.Id).First(&foundBot)
	logger.Infoln("Found bot: %v", foundBot)

	return foundBot
}

func (bot *Bot) Update(column string) bool {
	GetDB().Update(column, bot)
	return true
}
