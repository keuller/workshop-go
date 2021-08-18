package controller

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/keuller/account/internal/application"
	"github.com/keuller/account/internal/common"
)

type M map[string]string

type AccountController struct {
	accountService application.AccountService
}

func NewAccountController(svc application.AccountService) AccountController {
	return AccountController{svc}
}

func (ctrl AccountController) CreateAccountHandler(res http.ResponseWriter, req *http.Request) error {
	var data application.AccountRequest
	if err := common.BindJson(req.Body, &data); err != nil {
		return err
	}

	code, err := ctrl.accountService.CreateAccount(data)
	if err != nil {
		return err
	}

	common.ToJson(res, http.StatusCreated, M{
		"account_code": code,
		"message":      "Account was created successfully.",
	})

	return nil
}

func (ctrl AccountController) GetBalanceHandler(res http.ResponseWriter, req *http.Request) {
	account := chi.URLParam(req, "code")
	log.Printf("[INFO] get balance from account %s\n", account)
	result := ctrl.accountService.GetBalance(account)
	common.ToJson(res, http.StatusOK, result)
}

func (ctrl AccountController) Deposit(res http.ResponseWriter, req *http.Request) {
	var data application.DepositRequest
	if err := common.BindJson(req.Body, &data); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ctrl.accountService.Deposit(data); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	common.ToJson(res, http.StatusOK, M{
		"message": "Operation completed successfully.",
	})
}

func (ctrl AccountController) Transfer(res http.ResponseWriter, req *http.Request) {
	var data application.TransferRequest
	if err := common.BindJson(req.Body, &data); err != nil {
		log.Printf("[ERROR] fail to serialize JSON data, reason: %s \n", err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ctrl.accountService.Transfer(data); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	common.ToJson(res, http.StatusOK, M{
		"message": "Transfer OK",
	})
}
