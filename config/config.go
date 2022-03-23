package config

const (
	Api_key    string = "6212e60d29a7d50001efccd1"
	Passphrase string = "y9UU2JH9ZQUxgRjb1vHV8848DR1j17"
	Secret     string = "49abfb79-4d9b-4435-9ded-ab691e734d66"
)

const (
	DSize                   = "1"
	DSymbol                 = "SOL-USDT"
	DecimalPointNumber uint = 3
	DPriceMargin            = 0.005
	DChangeRate             = 0
)

var PriceMargin float64 = DPriceMargin //should be more than 0.002
var ChangeRate float64 = DChangeRate

func SetPriceMargin(newPrice float64) {
	PriceMargin = newPrice
}

func SetChangeRate(newPrice float64) {
	ChangeRate = newPrice
}
