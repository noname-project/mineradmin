package coin

func List() []string {
	return []string{
		"BTC",
		"BCH",
		"DASH",
		"ETH",
		"LTC",
	}
}

func Valid(coin string) bool {
	for _, c := range List() {
		if c == coin {
			return true
		}
	}
	return false
}
