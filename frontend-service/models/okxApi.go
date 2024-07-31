package models

import (
	"gorm.io/gorm"
	u "okx-bot/frontend-service/utils"
)

type OKxApi struct {
	gorm.Model
	MainKey    string `json:"mainKey"`
	SpecialKey string `json:"specialKey"`
	Phrase     string `json:"phrase"`
	UserId     uint   `json:"userId"`
}

const CREATED = "*** created ***"

func (api *OKxApi) Create() map[string]interface{} {
	GetDB().Create(api)

	if api.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	api.MainKey = CREATED
	api.SpecialKey = CREATED
	api.Phrase = CREATED

	response := u.Message(true, "Api has been created")
	response["api"] = api
	return response
}

func GetUserApi(userId uint) (*OKxApi, error) {
	api := &OKxApi{}
	err := GetDB().Table("o_kx_apis").Where("user_id = ?", userId).First(api).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return api, nil
}
