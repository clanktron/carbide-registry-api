package api

import (
	"carbide-registry-api/pkg/api/utils"
	DB "carbide-registry-api/pkg/database"
	"carbide-registry-api/pkg/objects"
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func releaseNameFromPath(r *http.Request) string {
	releaseName := r.PathValue("releaseName")
	return releaseName
}

// Responds with a JSON array of all releases in the database
//
// Success Code: 200 OK
func getAllReleasesHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var allReleases []objects.Release
		productName := productNameFromPath(r)
		limit, offset := utils.GetLimitAndOffset(r)
		allReleases, err := DB.GetAllReleasesforProduct(db, productName, limit, offset)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = utils.SendAsJSON(w, allReleases)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Accepts a JSON payload of a new release and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 Created
func createReleaseHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var newRelease objects.Release
		productName := productNameFromPath(r)
		err := utils.DecodeJSONObject(w, r, &newRelease)
		if err != nil {
			log.Error(err)
			return
		}
		parentProduct, err := DB.GetProduct(db, productName)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		newRelease.ProductId = &parentProduct.Id
		err = DB.AddRelease(db, newRelease)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		createdRelease, err := DB.GetRelease(db, newRelease)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"release": *createdRelease.Name,
		}).Info("new release has been successfully created")
		w.WriteHeader(http.StatusCreated)
		err = utils.SendAsJSON(w, createdRelease)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Responds with the JSON representation of a release
//
// Success Code: 200 OK
func getReleaseHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var retrievedRelease objects.Release
		productName := productNameFromPath(r)
		releaseName := releaseNameFromPath(r)
		retrievedRelease.Name = &releaseName
		parentProduct, err := DB.GetProduct(db, productName)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		retrievedRelease.ProductId = &parentProduct.Id
		retrievedRelease, err = DB.GetRelease(db, retrievedRelease)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = utils.SendAsJSON(w, retrievedRelease)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Accepts a JSON payload of the updated release and responds with the new JSON object after it's been successfully updated in the database
//
// Success Code: 200 OK
func updateReleaseHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var receivedRelease objects.Release
		productName := productNameFromPath(r)
		releaseName := releaseNameFromPath(r)
		err := utils.DecodeJSONObject(w, r, &receivedRelease)
		if err != nil {
			log.Error(err)
			return
		}
		receivedRelease.Name = &releaseName
		parentProduct, err := DB.GetProduct(db, productName)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		receivedRelease.ProductId = &parentProduct.Id
		err = DB.UpdateRelease(db, receivedRelease)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		updatedRelease, err := DB.GetRelease(db, receivedRelease)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"release": *updatedRelease.Name,
		}).Info("release has been successfully updated")
		w.WriteHeader(http.StatusOK)
		err = utils.SendAsJSON(w, updatedRelease)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Deletes the release and responds with an empty payload
//
// Success Code: 204 No Content
func deleteReleaseHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var releaseToDelete objects.Release
		productName := productNameFromPath(r)
		releaseName := releaseNameFromPath(r)
		releaseToDelete.Name = &releaseName
		parentProduct, err := DB.GetProduct(db, productName)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		releaseToDelete.ProductId = &parentProduct.Id
		err = DB.DeleteRelease(db, releaseToDelete)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"release": releaseName,
		}).Info("release has been successfully deleted")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	return http.HandlerFunc(fn)
}
