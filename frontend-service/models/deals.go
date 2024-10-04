package models

import (
	"gorm.io/gorm"
	"time"
)

type DealStatus int

const (
	DealStarted DealStatus = iota + 1
	DealFinished
	DealFailure
)

type DealDirection int

const (
	Long DealDirection = iota + 1
	Short
)

type Deal struct {
	gorm.Model
	StartTime   time.Time     `json:"startTime"`
	EndTime     time.Time     `json:"endTime"`
	StartAmount float64       `json:"startAmount"`
	EndAmount   float64       `json:"endAmount"`
	Status      DealStatus    `json:"status"`
	Direction   DealDirection `json:"direction"`
	BotRefer    uint          `json:"botRefer"`
	OrderId     string        `json:"orderId"`
	Stats       []DealStats   `json:"stats" gorm:"foreignKey:DealRefer"`
}

func (deal *Deal) StartDbSave(botId uint, startAmount float64) bool {
	var bot = Bot{ID: botId}
	db.Where("id = ?", botId).First(&bot)

	logger.Infoln("bot", bot)

	deal.StartAmount = startAmount
	deal.Status = DealStarted
	deal.StartTime = time.Now()

	bot.Deals = append(bot.Deals, *deal)
	db.Save(&bot)

	return true
}

func (deal *Deal) FinishDbSave(endAmount float64) bool {
	deal.EndAmount = endAmount
	deal.Status = DealFinished
	deal.EndTime = time.Now()
	db.Save(deal)

	return true
}

func FindByStatus(botId uint, status DealStatus) Deal {
	var foundDeal = Deal{}
	tx := db.Where("bot_refer = ? and status = ?", botId, status).First(&foundDeal)
	if tx.Error != nil {
		logger.Errorf("Found deal error: %v", tx.Error)
		return foundDeal
	}
	logger.Infoln("Found bot: %v", foundDeal)

	return foundDeal
}

func (deal *Deal) Failure(endAmount float64) bool {
	deal.EndAmount = endAmount
	deal.Status = DealFailure
	deal.EndTime = time.Now()
	db.Save(deal)

	return true
}

func (deal *Deal) Update(column string) bool {
	GetDB().Update(column, deal)
	return true
}
