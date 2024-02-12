package utils

import "net/http"

type RoundTrip func(req *http.Request) (*http.Response, error)

func (rt RoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}
