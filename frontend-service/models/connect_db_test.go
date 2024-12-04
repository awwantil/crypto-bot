package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	// Load test environment variables
	//err := godotenv.Load(".env.test")
	//assert.NoError(t, err)

	// Call ConnectDB function
	ConnectDB()

	// Check if db instance is created
	assert.NotNil(t, db)

	// Check if connection is successful
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())
}

//func TestGetDB(t *testing.T) {
//	// Set up test DB connection
//	dsn := "host=localhost user=testuser password=testpass dbname=testdb port=5432 sslmode=disable"
//	testDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	assert.NoError(t, err)
//	db = testDB
//
//	// Call GetDB function
//	retrievedDB := GetDB()
//
//	// Check if retrieved DB matches test DB
//	assert.Equal(t, testDB, retrievedDB)
//}

func TestGetDBNilConnection(t *testing.T) {
	// Set db to nil
	db = nil

	// Call GetDB function
	retrievedDB := GetDB()

	// Check if retrieved DB is nil
	assert.Nil(t, retrievedDB)
}
