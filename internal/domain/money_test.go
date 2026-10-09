package domain_test

import (
	"testing"
	"regexp"	

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


func TestParseMoneyMalformedInput(t *testing.T) {
	moneyPattern := regexp.MustCompile(`^-?[0-9]+\.[0-9]{2}$`)
	testCases := map[string]bool{
		"25.00":     true,
		"-25.00":    true,
		"0.00":      true,
		"":          false,
		"25":        false,
		"25.0":      false,
		"25.":       false,
		".25":       false,
		"--25.00":   false,
		"25..00":    false,
		"25.00.00":  false,
		" 25.00":    false,
		"+25.00":    false,
		"1e2":       false,
		"25.001":    false,
		"not-money": false,
	}

	for value, wantMatch := range testCases {
		t.Run(value, func(t *testing.T) {
			if got := moneyPattern.MatchString(value); got != wantMatch {
				t.Fatalf("money format match = %t, want %t", got, wantMatch)
			}

			_, err := domain.ParseMoney(value, "BRL")
			if wantMatch && err != nil {
				t.Fatalf("ParseMoney rejected valid value %q: %v", value, err)
			}
			if !wantMatch && err == nil {
				t.Fatalf("ParseMoney accepted malformed value %q", value)
			}
		})
	}
}

