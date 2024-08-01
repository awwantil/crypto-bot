package models

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var (
	blogger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.models.base",
	})
)
var db *gorm.DB

func ConnectDB() {
	blogger.Infoln("Connecting to DB")

	logger.Info("***** 5-5 *****")

	dbHost := os.Getenv("DB_HOST")
	logger.Info("dbHost: ", dbHost)

	if dbHost == "" {
		err := godotenv.Load()
		if err != nil {
			blogger.Errorf("Wrong load .env file: %s", err)
		}
	}

	tokenPassword := os.Getenv("TOKEN_PASSWORD")
	logger.Info("tokenPassword: ", tokenPassword)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	blogger.Infoln("dsn: ", dsn)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db = conn

	if err != nil {
		blogger.Errorf("Failed to connect to database: %s", err)
	} else {
		blogger.Infoln("connected")

		//blogger.Infoln("running migrations")
		//err = conn.Debug().AutoMigrate(&Account{}, &Contact{}, &TradingViewSignalReceive{}, &Signal{})
		//err = conn.Debug().AutoMigrate(&Bot{}, &OKxApi{}, &DealStats{}, &Deal{})
		//err = conn.Debug().AutoMigrate()
		if err != nil {
			blogger.Errorf("Failed to migrate database: %s", err)
		}
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		blogger.Error("Not db connection")
		return nil
	}
	return db
}
