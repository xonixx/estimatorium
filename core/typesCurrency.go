package core

import "strings"

type Currency uint8

const (
	CurrencyUnknown Currency = iota
	Usd
	Eur
)

var currency2Str = map[Currency]string{
	CurrencyUnknown: "CurrencyUnknown", Usd: "USD", Eur: "EUR",
}

func (c Currency) String() string {
	return currency2Str[c]
}

var currency2Symbol = map[Currency]string{
	Usd: "$", Eur: "€",
}

func (c Currency) Symbol() string {
	return currency2Symbol[c]
}

var currencyStr2Val = map[string]Currency{
	"USD": Usd, "EUR": Eur,
}

func CurrencyFromString(curr string) Currency {
	return currencyStr2Val[strings.ToUpper(curr)]
}
