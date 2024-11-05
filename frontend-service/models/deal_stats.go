package models

import (
	"gorm.io/gorm"
	"time"
)

type DealStats struct {
	gorm.Model
	CurrentTime time.Time `json:"startTime"`
	Value       float64   `json:"value"`
	DealRefer   uint      `json:"dealRefer"`
}
