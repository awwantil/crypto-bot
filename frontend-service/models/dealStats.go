package models

import (
	"gorm.io/gorm"
	"time"
)

type Stat struct {
	gorm.Model
	StatTime time.Time `json:"statTime"`
	Value    float64   `json:"value"`
}

type DealStats struct {
	gorm.Model
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Stats     []Stat    `json:"stats"`
}
