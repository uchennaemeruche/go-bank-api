package util

// const (
// 	USD = "USD"
// 	NGN = "NGN"
// 	EUR = "EUR"
// )

func getCurrencies() []string {
	return []string{"USD", "NGN", "EUR"}
}

func IsSupportedCurrencry(currency string) bool {
	for _, curr := range getCurrencies() {
		if curr == currency {
			return true
		}
	}
	return false
}
