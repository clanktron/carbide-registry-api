package api

import (
	"carbide-images-api/pkg/api/middleware"
	"carbide-images-api/pkg/api/utils"
	"carbide-images-api/pkg/objects"
	"net/http"
)

func authCheckHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !middleware.Authorized(r) {
			utils.HttpJSONError(w, "unauthorized", http.StatusUnauthorized)
		}
		return
	}
	return http.HandlerFunc(fn)
}

func loginHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var license objects.CarbideLicense
		if err := utils.DecodeJSONObject(w, r, &license); err != nil {
			return
		}
		if err := middleware.Login(w, license); err != nil {
			return
		}
	}
	return http.HandlerFunc(fn)
}
