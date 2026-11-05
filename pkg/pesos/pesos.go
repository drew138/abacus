package pesos

import (
	"math"
	"strconv"
	"strings"
)

func Format(v int64) string {
	sign := ""
	if v < 0 {
		sign, v = "-", -v
	}
	digits := strconv.FormatInt(v, 10)
	var b strings.Builder
	for i, d := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(d)
	}
	return sign + "$" + b.String()
}

func Short(v int64) string {
	switch {
	case v >= 1_000_000:
		s := strconv.FormatFloat(float64(v)/1_000_000, 'f', 1, 64)
		return "$" + strings.Replace(strings.TrimSuffix(s, ".0"), ".", ",", 1) + " M"
	case v >= 1_000:
		return "$" + strconv.FormatInt(int64(math.Round(float64(v)/1_000)), 10) + " k"
	}
	return "$" + strconv.FormatInt(v, 10)
}

func Percent(part, whole int64) string {
	if whole == 0 {
		return "0%"
	}
	return strings.Replace(strconv.FormatFloat(float64(part)*100/float64(whole), 'f', 1, 64), ".", ",", 1) + "%"
}
