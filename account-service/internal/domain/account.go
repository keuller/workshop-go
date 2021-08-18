package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string     `gorm:"size:36;primaryKey"`
	Owner     string     `gorm:"column:owner;size:36;not null"`
	Currency  string     `gorm:"column:currency;size:3;not null"`
	Balance   float64    `gorm:"column:balance"`
	CreatedAt time.Time  `gorm:"<-:create"`
	UpdatedAt *time.Time `gorm:"<-:update"`
}

type AccountBuilder struct {
	Account
}

func NewAccountBuilder() AccountBuilder {
	return AccountBuilder{}
}

// ====
func (b *AccountBuilder) WithCurrency(value string) *AccountBuilder {
	b.Currency = value
	return b
}

func (b *AccountBuilder) WithBalance(value float64) *AccountBuilder {
	b.Balance = value
	return b
}

func (b AccountBuilder) Build() Account {
	idValue := uuid.Must(uuid.NewRandom()).String()
	ownerId := uuid.Must(uuid.NewRandom()).String()

	if b.Owner != "" {
		ownerId = b.Owner
	}

	return Account{
		ID:        idValue,
		Owner:     ownerId,
		Balance:   b.Balance,
		Currency:  b.Currency,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}
}
