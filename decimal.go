package husky

import "github.com/shopspring/decimal"

func init() {
	decimal.DivisionPrecision = 18
}

func NewDeciamlForInt(v int64) decimal.Decimal {
	return decimal.NewFromInt(v)
}

func NewDeciamlForString(v string) (decimal.Decimal, error) {
	return decimal.NewFromString(v)
}
