package util

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const alphabet = "abcdefghijklmnopqrstuvqxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return cases.Title(language.Und).String(RandomString(6))
}

func RandomMoney() int64 {
	return RandomInt(0, 999999999)
}

func RandomCurrency() string {
	currencies := []string{"NGN", "USD", "EUR"}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}
