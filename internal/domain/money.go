package domain

// Money stores a currency amount as signed minor units.
type Money struct {
	amount   int64
	currency string
}

// ParseMoney parses a fixed-scale decimal amount into minor units.
func ParseMoney(value string, currency string) (Money, error) {
	return Money{}, nil
}

func (money Money) MinorUnits() int64 {
	return money.amount
}

func (money Money) Currency() string {
	return money.currency
}
