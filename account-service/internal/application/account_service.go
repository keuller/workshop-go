package application

import (
	"errors"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/keuller/account/internal/common"
	"github.com/keuller/account/internal/domain"
)

type AccountService struct {
	validate           *validator.Validate
	accountRepository  domain.IAccountRepository
	exchangeRepository domain.IExchangeRepository
}

func NewAccountService(repo domain.IAccountRepository, exchange domain.IExchangeRepository) AccountService {
	validate := validator.New()
	return AccountService{validate, repo, exchange}
}

// CreateAccount creates a new account
func (s AccountService) CreateAccount(data AccountRequest) (string, error) {
	if err := s.validate.Struct(data); err != nil {
		return "", common.BusinessFailure("failure on validation", err)
	}
	if !s.isValidCurrency(strings.ToLower(data.Currency)) {
		return "", common.BusinessFailure("Currency is invalid.", nil)
	}

	builder := domain.NewAccountBuilder()
	account := builder.WithBalance(0.0).
		WithCurrency(data.Currency).
		Build()

	if err := s.accountRepository.Create(account); err != nil {
		return "", err
	}

	return account.ID, nil
}

func (s AccountService) GetBalance(code string) BalanceResponse {
	account, err := s.accountRepository.FindByID(code)
	if err != nil {
		log.Printf("[FAIL] fail to get balance info from account, reason: %s", err.Error())
		return BalanceResponse{}
	}

	return BalanceResponse{
		Account:    code,
		Balance:    account.Balance,
		LastUpdate: common.DateToStr(account.CreatedAt),
	}
}

func (s AccountService) Deposit(data DepositRequest) error {
	account, err := s.accountRepository.FindByID(data.Account)
	if err != nil {
		log.Printf("[FAIL] account not found, reason: %s", err.Error())
		return errors.New("invalid account")
	}

	if err := s.updateBalance(account, data.Value, domain.CREDIT); err != nil {
		return err
	}

	return nil
}

func (s AccountService) Transfer(data TransferRequest) error {
	log.Printf("%v", data)
	if data.SourceAccount == data.TargetAccount {
		return errors.New("the accounts must be different")
	}

	source, err := s.accountRepository.FindByID(data.SourceAccount)
	if err != nil {
		log.Printf("[FAIL] source account not found, reason: %s", err.Error())
		return errors.New("invalid account")
	}

	target, err := s.accountRepository.FindByID(data.TargetAccount)
	if err != nil {
		log.Printf("[FAIL] target account not found, reason: %s", err.Error())
		return errors.New("invalid account")
	}

	if source.Balance < data.Value {
		return errors.New("insufficient balance on source account")
	}

	value := s.GetQuotationValue(source.Currency, target.Currency, data.Value)

	// debit from source account
	if err := s.updateBalance(source, data.Value, domain.DEBIT); err != nil {
		return err
	}

	// credit on target account
	if err := s.updateBalance(target, value, domain.CREDIT); err != nil {
		return err
	}

	log.Printf("Transfer from %s to %s, value: %.2f \n", source.Currency, target.Currency, value)
	return nil
}

func (s AccountService) isValidCurrency(value string) bool {
	for _, cur := range s.exchangeRepository.GetCurrencies() {
		if strings.ToLower(cur.Symbol) == value {
			return true
		}
	}
	return false
}

func (s AccountService) updateBalance(account domain.Account, value float64, operation int) error {
	newBalance := account.Balance + value
	if operation == domain.DEBIT {
		newBalance = account.Balance - value
	}

	if err := s.accountRepository.UpdateBalance(account.ID, newBalance); err != nil {
		return err
	}

	s.registerTransaction(account, operation, value)
	return nil
}

func (s AccountService) GetQuotationValue(sourceCurrency, targetCurrency string, value float64) float64 {
	if sourceCurrency == targetCurrency {
		return value
	}
	quotation := s.exchangeRepository.GetQuotation(sourceCurrency, targetCurrency, value)
	log.Printf("[DEBUG] quotation: %v \n", quotation)
	return quotation.GetValue()
}

func (s AccountService) registerTransaction(account domain.Account, operation int, value float64) {
	transactionBuilder := domain.NewTransactionBuilder()
	log.Printf("[DEBUG] registering transaction on account %s of operation %d with value %.2f", account.ID, operation, value)
	transaction := transactionBuilder.
		WithAccount(account).
		WithOperation(operation).
		WithValue(value).Build()
	s.accountRepository.AddTransaction(transaction)
}
