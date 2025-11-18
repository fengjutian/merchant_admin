package main

import (
	"log"
	"merchant_back/internal/config"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	db, err := cfg.InitDatabase()
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	DB = db
}

func CloseDB() {
	if err := config.CloseDatabase(DB); err != nil {
		log.Fatal("Failed to close database: ", err)
	}
	log.Println("Database connection closed")
}
