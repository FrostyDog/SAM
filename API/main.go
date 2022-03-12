package api

import "net/http"

var baseAPI string = "https://api.kucoin.com"

func Get(path string) {
	http.Get(baseAPI + path)
}
