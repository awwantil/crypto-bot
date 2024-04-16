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
	logger.Info("4 - !!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	logger.Infof("dbName: %s, dbPassword: %s\n", dbName, dbPassword)

	hostName := os.Getenv("DB_HOST")
	if os.Getenv("NODE_ENV") == "dev" {
		hostName = "localhost"
	} else {
		blogger.Infoln("Waiting 5 second ...")
		//time.Sleep(5 * time.Second)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		hostName,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	blogger.Infoln("dsn: ", dsn)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db = conn

	if err != nil {
		blogger.Errorf("Failed to connect to database: %s", err)
	}

	blogger.Infoln("connected")

	blogger.Infoln("running migrations")
	err = conn.Debug().AutoMigrate(&Account{}, &Contact{}, &TradingViewSignalReceive{}, &Signal{})
	err = conn.Debug().AutoMigrate(&Bot{})
	if err != nil {
		blogger.Errorf("Failed to migrate database: %s", err)
	}

}

func GetDB() *gorm.DB {
	if db == nil {
		blogger.Error("Not db connection")
		return nil
	}
	return db
}
