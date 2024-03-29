package repository

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/keuller/account/internal/common"
	"github.com/keuller/account/internal/domain"
	"github.com/keuller/account/internal/infra"
)

const (
	timeout = 1500 * time.Millisecond
)

type exchangeRepository struct {
	exchangeUrl string
	dataCache   []domain.Currency
}

func NewExchangeRepository() domain.IExchangeRepository {
	exchangeUrl := infra.GetConfig().ExchangeUrl
	cached := make([]domain.Currency, 0)
	return &exchangeRepository{exchangeUrl, cached}
}

func (er *exchangeRepository) GetCurrencies() []domain.Currency {
	if len(er.dataCache) > 0 {
		return er.dataCache
	}

	currencies := make([]domain.Currency, 0)
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(er.exchangeUrl + "/v1/exchange/currencies")
	if err != nil {
		log.Printf("[FAIL] fail to get currencies, reason: %s.\n", err.Error())
		return currencies
	}

	if resp.StatusCode != 200 {
		log.Printf("[FAIL] exchange service returned invalid code: %d \n", resp.StatusCode)
		return currencies
	}

	if err := common.BindJson(resp.Body, &currencies); err != nil {
		log.Printf("[FAIL] cannot decode message from exchange service: %s \n", err.Error())
	}

	er.dataCache = currencies
	return currencies
}

func (er exchangeRepository) GetQuotation(from, to string, value float64) domain.Quotation {
	var quotation domain.Quotation
	client := &http.Client{
		Timeout: timeout,
	}

	url := fmt.Sprintf("%s/v1/exchange/quotation?from=%s&to=%s&val=%.3f", er.exchangeUrl, from, to, value)
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("[FAIL] fail to get quotation, reason: %s.\n", err.Error())
		return quotation
	}

	if resp.StatusCode != 200 {
		log.Printf("[FAIL] exchange service returned invalid code: %d \n", resp.StatusCode)
		return domain.Quotation{}
	}

	if err := common.BindJson(resp.Body, &quotation); err != nil {
		log.Printf("[FAIL] cannot decode message from exchange service: %s \n", err.Error())
	}

	return quotation
}
