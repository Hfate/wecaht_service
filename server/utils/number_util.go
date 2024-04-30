package utils

import (
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"math"
	"strings"
)

func Float64MulRound2(num1 float64, num2 float64) float64 {
	result, _ := decimal.NewFromFloat(num1).Mul(decimal.NewFromFloat(num2)).Round(2).Float64()
	return result
}

func Float64Fixed(num float64, fixed int) string {
	decimalNum := decimal.NewFromFloat(num)
	fixedNum := decimalNum.StringFixed(int32(fixed))
	return fixedNum
}

func StringFixed2(num string) string {
	return StringFixed(num, 2)
}

func StringFixed(num string, fixed int) string {
	decimalNum, _ := decimal.NewFromString(num)
	fixedNum := decimalNum.StringFixed(int32(fixed))
	return fixedNum
}

func FloatCmp(num1 string, num2 string) int {
	num1D, _ := decimal.NewFromString(num1)
	num2D, _ := decimal.NewFromString(num2)
	result := num1D.Cmp(num2D)
	return result
}

func FloatAdd(num1 string, num2 string) string {
	num1D, _ := decimal.NewFromString(num1)
	num2D, _ := decimal.NewFromString(num2)
	result := num1D.Add(num2D).String()
	return result
}

func FloatSub(num1 string, num2 string) string {
	num1D, _ := decimal.NewFromString(num1)
	num2D, _ := decimal.NewFromString(num2)

	result := num1D.Sub(num2D).String()

	return result
}

func FloatMul(num1 string, num2 string) string {
	num1D, _ := decimal.NewFromString(num1)
	num2D, _ := decimal.NewFromString(num2)

	result := num1D.Mul(num2D).String()

	return result
}

func FloatEqual(num1 string, num2 string) bool {

	num1D, _ := decimal.NewFromString(num1)
	num2D, _ := decimal.NewFromString(num2)

	return num1D.Equal(num2D)
}

func FloatEqualZero(num1 string) bool {
	num1D, _ := decimal.NewFromString(num1)
	num2D, _ := decimal.NewFromString("0")

	return num1D.Equal(num2D)
}

func FloatNotEqualZero(num1 string) bool {
	return !FloatEqualZero(num1)
}

func FloatDiv(num1 string, num2 string) string {

	if num2 == "0" || num1 == "0" {
		return "0"
	}
	// 被除数为0
	if cast.ToFloat64(num2) == 0 {
		return "0"
	}

	num1D, _ := decimal.NewFromString(num1)
	num2D, _ := decimal.NewFromString(num2)

	result := num1D.Div(num2D).String()

	return result
}

func FloatRoundTen(num1 string) string {
	return StringFixed(num1, 10)
}

func FloatRoundUp(original string, roundNumber int) string {

	num1D, _ := decimal.NewFromString(original)
	num2D := decimal.NewFromInt(int64(roundNumber))

	return num1D.Div(num2D).Ceil().Mul(num2D).String()
}

func FloatRound(num1 string, round int) float64 {
	num1DStr, _ := decimal.NewFromString(num1)
	result, _ := num1DStr.Round(int32(round)).Float64()
	return result
}

func FloatCeil(num1 string) string {
	num1DStr, _ := decimal.NewFromString(num1)
	result := num1DStr.Ceil().String()
	return result
}

func Float64Round2Str(num float64) string {
	return Float64Fixed(num, 2)
}

func Float64Round2(num float64) float64 {
	result, _ := decimal.NewFromFloat(num).Round(int32(2)).Float64()
	return result
}

func AbsNum(num float64) float64 {
	return math.Abs(num)
}

func Float64DivideRound2(num1 float64, num2 float64) float64 {
	if num2 == 0 || num1 == 0 {
		return 0
	}
	result, _ := decimal.NewFromFloat(num1).Div(decimal.NewFromFloat(num2)).Round(2).Float64()
	return result
}

func IntDivideRound2(num1 int, num2 int) float64 {
	return Float64DivideRound2(float64(num1), float64(num2))
}

func Float64DivideRound1(num1 float64, num2 float64) float64 {
	if num2 == 0 || num1 == 0 {
		return 0
	}
	result, _ := decimal.NewFromFloat(num1).Div(decimal.NewFromFloat(num2)).Round(1).Float64()
	return result
}

func StrArr2Uint64Arr(nums []string) []uint64 {
	arr := make([]uint64, 0)
	for _, s := range nums {
		if s == "" {
			continue
		}
		arr = append(arr, cast.ToUint64(s))
	}
	return arr
}

func Uint64Arr2StrArr(nums []uint64) []string {
	arr := make([]string, 0)
	for _, s := range nums {
		arr = append(arr, cast.ToString(s))
	}
	return arr
}

func ParseMoney(money string) string {
	if money == "" {
		return money
	}

	return strings.ReplaceAll(money, ",", "")
}
