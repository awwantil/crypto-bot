package models

import (
	"fmt"
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
	UserId        uint      `json:"userId"`
	StartTime     time.Time `json:"startTime"`
	Status        BotStatus `json:"status"`
	InitialAmount float64   `json:"initialAmount"`
	CurrentAmount float64   `json:"currentAmount"`
	Lever         float64   `json:"lever"`
	PosSide       uint      `json:"posSide"`
	SignalRefer   uint      `json:"signalRefer"`
	OkxSignalId   string    `json:"okxSignalId"`
	OkxBotId      string    `json:"okxBotId"`
	DealsPercent  float64   `json:"dealsPercent"`
	IsProduction  bool      `json:"isProduction"`
	Deals         []Deal    `json:"deals" gorm:"foreignKey:BotRefer"`
}

type BotCreateRequest struct {
	InitialAmount float64 `json:"initialAmount"`
	CodeSignalId  string  `json:"codeSignalId"`
	Lever         float64 `json:"lever"`
	DealsPercent  float64 `json:"dealsPercent"`
	IsProduction  bool    `json:"isProduction"`
}

type BotWithIdRequest struct {
	Id uint `json:"id"`
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

func (bot *Bot) Create(signal *Signal) map[string]interface{} {

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

func FindBotByByOkxSignalId(signalId string) (*Bot, error) {
	findBot := &Bot{}
	err := GetDB().Table("bots").Where("okx_signal_id = ?", signalId).Find(&findBot).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return findBot, nil
}

func (bot *Bot) Update() bool {
	GetDB().Save(bot)
	return true
}
