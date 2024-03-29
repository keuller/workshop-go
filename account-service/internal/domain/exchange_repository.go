package domain

import "strconv"

type Currency struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type Quotation struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func (q Quotation) GetValue() float64 {
	val, _ := strconv.ParseFloat(q.Value, 64)
	return val
}

type IExchangeRepository interface {
	GetCurrencies() []Currency
	GetQuotation(from, to string, value float64) Quotation
}
