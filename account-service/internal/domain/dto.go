package domain

import "fmt"

type AccountRequest struct {
	Name     string `json:"name" validate:"required,min=5"`
	Email    string `json:"email" validate:"required,email"`
	Document string `json:"document" validate:"required,min=7"`
	Currency string `json:"currency" validate:"required,len=3"`
	Age      int    `json:"age" validate:"required,min=16"`
}

type BalanceResponse struct {
	Account    string  `json:"account"`
	LastUpdate string  `json:"last_update"`
	Balance    float64 `json:"balance"`
}

type DepositRequest struct {
	Account string  `json:"account_code" validate:"required,len=36"`
	Value   float64 `json:"value" validate:"required,gte=0"`
}

type TransferRequest struct {
	SourceAccount string  `json:"from"`
	TargetAccount string  `json:"to"`
	Value         float64 `json:"value"`
}

func (tr TransferRequest) String() string {
	return fmt.Sprintf("[source: %s, target: %s, value: %.2f", tr.SourceAccount, tr.TargetAccount, tr.Value)
}
