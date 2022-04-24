package config

const (
	Api_key    string = "6212e60d29a7d50001efccd1"
	Passphrase string = "y9UU2JH9ZQUxgRjb1vHV8848DR1j17"
	Secret     string = "49abfb79-4d9b-4435-9ded-ab691e734d66"
)

const (
	TradingSize             = "1.5"
	PrimarySymbol           = "SOL"
	SecondarySymbol         = "USDT"
	DecimalPointNumber uint = 3
	BaseMargin              = 0.007 // more than 0.002
	DChangeRate             = 0
)

var TradingPair string = PrimarySymbol + "-" + SecondarySymbol
var ChangeRate float64 = DChangeRate

func SetChangeRate(newPrice float64) {
	ChangeRate = newPrice
}
