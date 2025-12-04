package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"vibe-music-server/internal/config"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Init() {
	once.Do(func() {
		DB := config.Get().Database
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			DB.Username, DB.Password, DB.Host, DB.Name)
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}
	})
}

func Get() *gorm.DB {
	return db
}
