package kucoin_api

import (
	"github.com/FrostyDog/SAM/config"
	"github.com/Kucoin/kucoin-go-sdk"
)

var S *kucoin.ApiService = kucoin.NewApiService(
	kucoin.ApiBaseURIOption("https://api.kucoin.com"),
	kucoin.ApiKeyOption(config.Api_key),
	kucoin.ApiSecretOption(config.Secret),
	kucoin.ApiPassPhraseOption(config.Passphrase),
	kucoin.ApiKeyVersionOption(kucoin.ApiKeyVersionV2))
