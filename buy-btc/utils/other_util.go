package utils

import (
	"math"
)

// 少数切り捨て
func RoundDecimal(num float64) float64 {
	return math.Round(num)
}

// 少数切り上げ
func roundUp(num, places float64) float64 {
	shift := math.Pow(10, places)
	return RoundDecimal(num*shift) / shift
}

// 予算を元に購入数量を算出する
func CalcAmount(price, budget, minmumAmount, places float64) float64 {
	amount := roundUp(budget/price, places)
	if amount < minmumAmount {
		return minmumAmount
	}
	return amount
}
