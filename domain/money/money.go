package money

type Value int

func (v Value) Int() int {
	return int(v)
}

type Currency string

func (c Currency) String() string {
	return string(c)
}

type Money struct {
	Value    Value    `json:"value"`
	Currency Currency `json:"currency"`
}

func NewMoney(value int, currency string) Money {
	return Money{
		Value:    Value(value),
		Currency: Currency(currency),
	}
}
