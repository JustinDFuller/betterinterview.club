package organization

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./organization/notfound.html")
}
