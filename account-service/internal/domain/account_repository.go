package domain

type IAccountRepository interface {
	Create(data Account) error
	FindByID(id string) (Account, error)
	UpdateBalance(id string, value float64) error
	AddTransaction(data Transaction)
}
