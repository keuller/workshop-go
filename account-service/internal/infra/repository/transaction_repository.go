package repository

import (
	"gorm.io/gorm"

	"github.com/keuller/account/internal/domain"
)

type transactionRepository struct {
}

func (r transactionRepository) Add(db *gorm.DB, data domain.Transaction) error {
	if res := db.Table("transactions").Create(data); res.Error != nil {
		return res.Error
	}
	return nil
}
