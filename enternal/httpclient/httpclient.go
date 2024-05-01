package httpclient

import "net/http"

func NewHttpClient() *http.Client {
	c := http.Client{
		Transport:     NewLogginRoundTripper(),
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	return &c
}
