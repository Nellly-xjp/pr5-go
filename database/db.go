package database

import (
	"log"
	"pr-5-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=12345678incidents_db port=5432 sslmode=disable TimeZone=Europe/Kyiv"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Помилка підключення до PostgreSQL: ", err)
	}

	err = DB.AutoMigrate(&models.Incident{}, &models.User{})
	if err != nil {
		log.Fatal("Помилка авто-міграції БД: ", err)
	}

	log.Println("Успішне підключення до PostgreSQL та міграція таблиць!")
}