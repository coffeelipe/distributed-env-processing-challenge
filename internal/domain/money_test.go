package domain_test

import (
	"testing"

	"coffeelipe/distributed-processing/internal/domain"
)

func TestParseMoney(t *testing.T) {
	money, err := domain.ParseMoney("25.00", "BRL")
	if err != nil {
		t.Fatalf("parse money: %v", err)
	}

	if got := money.MinorUnits(); got != 2500 {
		t.Fatalf("minor units = %d, want 2500", got)
	}

	if got := money.Currency(); got != "BRL" {
		t.Fatalf("currency = %q, want BRL", got)
	}
}
