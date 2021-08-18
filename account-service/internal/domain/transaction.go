package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	INVALID = iota
	CREDIT
	DEBIT
)

type Transaction struct {
	ID        string    `gorm:"size:36;primaryKey"`
	Account   string    `gorm:"column:account;size:36;not null"`
	Currency  string    `gorm:"column:currency;size:3;not null"`
	Operation int       `gorm:"column:operation"`
	Value     float64   `gorm:"column:value"`
	CreatedAt time.Time `gorm:"<-:create"`
}

type TransactionBuilder struct {
	Transaction
}

func NewTransactionBuilder() TransactionBuilder {
	return TransactionBuilder{}
}

// -------------------------------------- \\
func (b *TransactionBuilder) WithAccount(acc Account) *TransactionBuilder {
	b.Account = acc.ID
	b.Currency = acc.Currency
	return b
}

func (b *TransactionBuilder) WithOperation(value int) *TransactionBuilder {
	b.Operation = value
	return b
}

func (b *TransactionBuilder) WithValue(value float64) *TransactionBuilder {
	b.Value = value
	return b
}

func (b TransactionBuilder) Build() Transaction {
	idValue := uuid.Must(uuid.NewRandom())
	return Transaction{
		ID:        idValue.String(),
		Account:   b.Account,
		Currency:  b.Currency,
		Operation: b.Operation,
		Value:     b.Value,
		CreatedAt: time.Now(),
	}
}
