package config

type BotGlobalConfig struct {
	Name         string
	ExchangeName string
}

type CEXConfigs struct {
	ExchangeName string
	Api_key      string
	Passphrase   string
	Secret       string
}

var KucoinConfig CEXConfigs = CEXConfigs{
	ExchangeName: "KuCoin",
	Api_key:      "6212e60d29a7d50001efccd1",
	Passphrase:   "y9UU2JH9ZQUxgRjb1vHV8848DR1j17",
	Secret:       "49abfb79-4d9b-4435-9ded-ab691e734d66"}
