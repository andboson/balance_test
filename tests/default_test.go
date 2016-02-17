package test

import (
	"app/services"
	"app/models"
)

func init() {
	services.InitDb()
	services.DB.AutoMigrate(&models.UserBalance{})
}

///mocks
func AddMockDataModel(id, balance int) {
	services.DB.Where("id = ?", id).Delete(models.UserBalance{})
	models.CreateUserBalance(id, balance, "")
}

func DeleteMockDataModel(id int) {
	services.DB.Where("id = ?", id).Delete(models.UserBalance{})
}
