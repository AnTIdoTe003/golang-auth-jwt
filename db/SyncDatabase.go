package db

import "example/REST_API_JWT/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}