package models

import (
	"gorm.io/gorm"
	"time"
)

type Deal struct {
	gorm.Model
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
