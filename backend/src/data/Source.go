package utilities

import (
	"backend/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"os"
)

func SetupDB() *gorm.DB {
	godotenv.Load(".env")
	USER := os.Getenv("DB_USERNAME")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	PORT := os.Getenv("DB_PORT")
	DBNAME := os.Getenv("DB_DATABASE")
	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	db, err := gorm.Open(os.Getenv("DB_CONNECTION"), URL)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Config() {
	db := SetupDB()
}
