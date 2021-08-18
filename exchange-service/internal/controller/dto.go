package controller

type CurrencyResponse struct {
	Symbol string `json:"symbol" msg:"symbol"`
	Name   string `json:"name" msg:"name"`
}

type QuotationResponse struct {
	From  string `json:"from" codec:"from"`
	To    string `json:"to" codec:"to"`
	Value string `json:"value" codec:"value"`
}
