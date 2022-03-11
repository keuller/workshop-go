package application

import (
	"testing"

	"github.com/keuller/account/internal/domain"
	"github.com/keuller/account/test/mock"
)

func TestAccountServiceSuite(t *testing.T) {
	accountRepo := mock.NewAccountRepositoryMock(false)
	exchangeRepo := mock.NewExchangeRepositoryMock()

	t.Run("Create a valid account", func(it *testing.T) {
		service := domain.NewAccountService(accountRepo, exchangeRepo)
		data := domain.AccountRequest{
			Name:     "John Doe",
			Document: "8546320",
			Email:    "john.doe@outlook.com",
			Currency: "USD",
			Age:      26,
		}

		code, err := service.CreateAccount(data)
		if err != nil {
			it.Fatalf("%v", err)
		}
		if code == "" {
			it.Fatal("Expected code not empty")
		}
	})

	t.Run("Try to create an account with invalid age", func(it *testing.T) {
		service := domain.NewAccountService(accountRepo, exchangeRepo)
		data := domain.AccountRequest{
			Name:     "John Doe",
			Document: "8546320",
			Email:    "john.doe@outlook.com",
			Currency: "USD",
			Age:      14,
		}

		_, err := service.CreateAccount(data)
		if err == nil {
			it.Fatal("Expected error for invalid age.")
		}

		// validationErrors := err.(validator.ValidationErrors)
		// for _, errs := range validationErrors {
		// 	fmt.Println(errs)
		// }
	})

	t.Run("Try to create an account with invalid currency", func(it *testing.T) {
		service := domain.NewAccountService(accountRepo, exchangeRepo)
		data := domain.AccountRequest{
			Name:     "John Doe",
			Email:    "john.doe@outlook.com",
			Currency: "CAD",
			Age:      18,
		}

		_, err := service.CreateAccount(data)
		if err == nil {
			it.Fatal("Expected error for invalid currency.")
		}
	})

	t.Run("Get Balance", func(it *testing.T) {
		service := domain.NewAccountService(accountRepo, exchangeRepo)
		response := service.GetBalance("8695230152")
		if response.Balance != 1.0 && response.Account != "73089a35-3b88-40df-ab27-64df5e58e343" {
			it.Fatal("failed to get balence from account")
		}
	})

	t.Run("Account deposit", func(it *testing.T) {
		service := domain.NewAccountService(accountRepo, exchangeRepo)
		data := domain.DepositRequest{
			Account: "8695230152",
			Value:   10.0,
		}
		if err := service.Deposit(data); err != nil {
			it.Fatalf("Error non expected, but %v", err)
		}
	})
}
