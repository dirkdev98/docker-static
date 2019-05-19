package static

import (
	"net/http"
	"strconv"
)

var (
	cacheControl = new(header)
)

type header struct {
	name  string
	value string
}

func addCacheControlValue(maxAge int) {
	cacheControl.name = "Cache-Control"
	cacheControl.value = "max-age=" + strconv.FormatInt(int64(maxAge), 10)
}

func addHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(cacheControl.name, cacheControl.value)

		next(w, r)
	}
}
