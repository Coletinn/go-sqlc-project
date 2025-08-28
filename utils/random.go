package utils

import "github.com/brianvoe/gofakeit/v7"

func RandomOwner() string {
	return gofakeit.Name()
}

func RandomMoney() float64 {
	return gofakeit.Price(10, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "JPY"}
	return currencies[gofakeit.Number(0, len(currencies)-1)]
}
