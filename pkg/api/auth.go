package api

import (
	"carbide-images-api/pkg/api/middleware"
	"carbide-images-api/pkg/api/utils"
	DB "carbide-images-api/pkg/database"
	"carbide-images-api/pkg/objects"
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func authCheckHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		middleware.Authorized(w, r)
		return
	}
	return http.HandlerFunc(fn)
}

func loginHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var userLoggingIn objects.User
		err := utils.DecodeJSONObject(w, r, &userLoggingIn)
		if err != nil {
			log.Error(err)
			return
		}
		if err := DB.VerifyUser(db, userLoggingIn); err != nil {
			log.Error(err)
			utils.HttpJSONError(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		userLoggingIn, err = DB.GetUser(db, *userLoggingIn.Username)
		if err != nil {
			log.Error(err)
			return
		}
		err = middleware.Authenticate(w, userLoggingIn)
		if err != nil {
			log.Error(err)
		}
		log.WithFields(log.Fields{
			"user": *userLoggingIn.Username,
		}).Info("user logged in successfully")
		utils.RespondWithJSON(w, "login successfull")
	}
	return http.HandlerFunc(fn)
}
