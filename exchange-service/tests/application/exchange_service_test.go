package application

import (
	"fmt"
	"testing"

	"github.com/keuller/exchange/internal/application"
)

func TestExchangeService(t *testing.T) {

	t.Run("Conversion from USD to BRL", func(it *testing.T) {
		quotationSvc := application.NewExchangeService()
		quotation, _ := quotationSvc.Get("USD", "BRL", 10.0)
		result := fmt.Sprintf("%.2f", quotation.Result)
		if result != "50.74" {
			it.Fatalf("Expected: 50.74, Found: %s", result)
		}
	})
}
