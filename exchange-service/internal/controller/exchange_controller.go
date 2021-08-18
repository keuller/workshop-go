package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/keuller/exchange/internal/application"
	"github.com/keuller/exchange/internal/domain"
	"github.com/ugorji/go/codec"
)

var (
	msgHandle codec.MsgpackHandle
)

func GetCurrencies(res http.ResponseWriter, req *http.Request) {
	currency := domain.NewCurrency()
	currencies := make([]CurrencyResponse, 0)
	for _, item := range currency.GetAll() {
		currencies = append(currencies, CurrencyResponse{item.Symbol, item.Name})
	}

	Json(res, currencies)
}

func GetQuotation(res http.ResponseWriter, req *http.Request) {
	from := strings.ToUpper(req.URL.Query().Get("from"))
	to := strings.ToUpper(req.URL.Query().Get("to"))
	value, err := strconv.ParseFloat(req.URL.Query().Get("val"), 32)
	if sendError(res, err, http.StatusBadRequest) {
		return
	}

	quotationSvc := application.NewExchangeService()
	quotation, err := quotationSvc.Get(from, to, value)
	if sendError(res, err, http.StatusBadRequest) {
		return
	}

	log.Printf("[INFO] get quotation from %s to %s with %.2f", from, to, value)
	response := QuotationResponse{
		From:  from,
		To:    to,
		Value: fmt.Sprintf("%.2f", quotation.Result),
	}

	Json(res, response)
}

func GetQuotationMsg(res http.ResponseWriter, req *http.Request) {
	from := strings.ToUpper(req.URL.Query().Get("from"))
	to := strings.ToUpper(req.URL.Query().Get("to"))
	value, err := strconv.ParseFloat(req.URL.Query().Get("val"), 32)
	if sendError(res, err, http.StatusBadRequest) {
		return
	}

	quotationSvc := application.NewExchangeService()
	quotation, err := quotationSvc.Get(from, to, value)
	if sendError(res, err, http.StatusBadRequest) {
		return
	}

	log.Printf("[INFO] get quotation from %s to %s with %.2f", from, to, value)
	response := QuotationResponse{
		From:  from,
		To:    to,
		Value: fmt.Sprintf("%.2f", quotation.Result),
	}

	buf := &bytes.Buffer{}
	enc := codec.NewEncoder(buf, &msgHandle)

	if err := enc.Encode(response); sendError(res, err, http.StatusInternalServerError) {
		return
	}

	res.Header().Set("Content-Type", "application/octet-stream")
	res.WriteHeader(http.StatusOK)
	res.Write(buf.Bytes())
}

// Json formats the response to 'application/json' type
func Json(res http.ResponseWriter, data interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(data); sendError(res, err, http.StatusInternalServerError) {
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write(buf.Bytes())
}

func sendError(res http.ResponseWriter, err error, code int) bool {
	if err == nil {
		return false
	}

	http.Error(res, err.Error(), code)
	return true
}
