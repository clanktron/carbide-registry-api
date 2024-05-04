package middleware

import (
	"net/http"
)

func ChainHandlers(handlers ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
    return func(finalHandler http.Handler) http.Handler {
		for _, handler := range handlers {
            finalHandler = handler(finalHandler)
        }
        return finalHandler
    }
}

func Healthcheck() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		return
	}
	return http.HandlerFunc(fn)
}

func CORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		enableCors(w, r)
		if r.Method == "OPTIONS" {
			return
		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
