package kucoin_api

import (
	"github.com/FrostyDog/SAM/config"
	"github.com/Kucoin/kucoin-go-sdk"
)

var S *kucoin.ApiService = kucoin.NewApiService(
	kucoin.ApiBaseURIOption("https://api.kucoin.com"),
	kucoin.ApiKeyOption(config.KucoinConfig.Api_key),
	kucoin.ApiSecretOption(config.KucoinConfig.Secret),
	kucoin.ApiPassPhraseOption(config.KucoinConfig.Passphrase),
	kucoin.ApiKeyVersionOption(kucoin.ApiKeyVersionV2))
