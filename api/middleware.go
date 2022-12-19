package api

import "net/http"

func middleware(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return withGzip(withContentSecurityPolicy(handler))
}
