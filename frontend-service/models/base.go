package models

import (
	"fmt"
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

	//e := godotenv.Load()
	//if e != nil {
	//	fmt.Print(e)
	//}

	logger.Info("***** 1 *****")

	dbHost := os.Getenv("DB_HOST")
	logger.Info("dbHost", dbHost)

	dbUser := os.Getenv("DB_USER")
	logger.Info("dbUser", dbUser)

	dbPassword := os.Getenv("DB_PASSWORD")
	logger.Info("dbPassword", dbPassword)

	dbName := os.Getenv("DB_NAME")
	logger.Info("dbName", dbName)

	tokenPassword := os.Getenv("TOKEN_PASSWORD")
	logger.Info("tokenPassword", tokenPassword)

	nodeEnv := os.Getenv("NODE_ENV")
	logger.Info("nodeEnv", nodeEnv)

	dsn := fmt.Sprintf(
		"host=185.240.103.204 user=%s password=%s dbname=%s port=5432 sslmode=disable",
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

		blogger.Infoln("running migrations")
		err = conn.Debug().AutoMigrate(&Account{}, &Contact{}, &TradingViewSignalReceive{}, &Signal{})
		err = conn.Debug().AutoMigrate(&Bot{})
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
