package domain_test

import (
	"regexp"
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

func TestParseMoneyInt64Bounds(t *testing.T) {
	const (
		maxMinorUnits int64 = 1<<63 - 1
		minMinorUnits int64 = -1 << 63
	)

	testCases := []struct {
		name      string
		value     string
		wantMinor int64
		wantError bool
	}{
		{name: "overflow one minor unit above maximum", value: "92233720368547758.08", wantError: true},
		{name: "underflow one minor unit below minimum", value: "-92233720368547758.09", wantError: true},
		{name: "grossly above maximum", value: "9223372036854775807.00", wantError: true},
		{name: "grossly below minimum", value: "-9223372036854775808.00", wantError: true},
		{name: "exact maximum", value: "92233720368547758.07", wantMinor: maxMinorUnits},
		{name: "exact minimum", value: "-92233720368547758.08", wantMinor: minMinorUnits},
		{name: "zero", value: "0.00", wantMinor: 0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			money, err := domain.ParseMoney(testCase.value, "BRL")
			if testCase.wantError {
				if err == nil {
					t.Fatalf("ParseMoney accepted out-of-bounds value %q", testCase.value)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseMoney rejected in-bounds value %q: %v", testCase.value, err)
			}
			if got := money.MinorUnits(); got != testCase.wantMinor {
				t.Fatalf("minor units = %d, want %d", got, testCase.wantMinor)
			}
		})
	}
}
