package repository

import (
	"log"

	"gorm.io/gorm"

	"github.com/keuller/account/internal/domain"
)

type accountRepository struct {
	DB *gorm.DB
	transactionRepository
}

func NewAccountRepository(conn *gorm.DB) domain.IAccountRepository {
	return accountRepository{DB: conn}
}

func (r accountRepository) Create(data domain.Account) error {
	if res := r.DB.Create(data); res.Error != nil {
		return res.Error
	}
	return nil
}

func (r accountRepository) FindByID(id string) (domain.Account, error) {
	var account domain.Account
	res := r.DB.Table("accounts").Where("id = ?", id).First(&account)
	if res.Error != nil {
		return domain.Account{}, res.Error
	}
	return account, nil
}

func (r accountRepository) UpdateBalance(id string, value float64) error {
	res := r.DB.Table("accounts").
		Where("id = ?", id).
		Update("balance", value)
	
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r accountRepository) AddTransaction(data domain.Transaction) {
	if err := r.transactionRepository.Add(r.DB, data); err != nil {
		log.Printf("[FAIL] cannot generate transaction, reason: %s \n", err.Error())
	}
}
