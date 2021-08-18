package domain

type ITransactionRepository interface {
	Create(data Transaction) error
}
