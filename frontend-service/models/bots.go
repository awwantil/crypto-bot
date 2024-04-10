package models

import (
	"gorm.io/gorm"
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
	StartTime     time.Time    `json:"startTime"`
	Status        BotStatus    `json:"status"`
	InitialAmount float64      `json:"initialAmount"`
	CurrentAmount float64      `json:"currentAmount"`
	Signal        SignalObject `json:"signal"`
	//Deals     []Deal       `json:"deals" gorm:"foreignKey:UserRefer"`
}

func (index BotStatus) String() string {
	return [...]string{"Created", "Waiting", "MakingDeal", "Stopped", "Error", "Pause"}[index-1]
}

func (index BotStatus) EnumIndex() int {
	return int(index)
}
