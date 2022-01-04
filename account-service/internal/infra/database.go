package infra

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/keuller/account/internal/domain"
)

var (
	_conn *gorm.DB
)

func InitDB() error {
	dbFile := "./" + GetConfig().DbFile
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return err
	}
	_conn = db

	// execute the migrations
	db.AutoMigrate(&domain.Account{})
	db.AutoMigrate(&domain.Transaction{})

	return nil
}

func Connection() *gorm.DB {
	return _conn
}
