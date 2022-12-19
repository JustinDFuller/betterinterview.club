package api

import "net/http"

func withContentSecurityPolicy(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'none'; connect-src 'self'; img-src 'self'; style-src 'unsafe-inline'; form-action 'self'; manifest-src 'self'; script-src 'sha256-g1OdTc+cgep4LHYKLdcx6nMjok9omWedJvE365v/NGE='; worker-src 'self'; report-uri /csp-violation/")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("Referrer-Policy", "same-origin")
		handler(w, r)
	}
}
