package api

import (
	"carbide-registry-api/pkg/api/utils"
	DB "carbide-registry-api/pkg/database"
	"carbide-registry-api/pkg/objects"
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Responds with a JSON array of all images in the database
//
// Success Code: 200 OK
func getAllReleaseImageMappingsHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		limit, offset := utils.GetLimitAndOffset(r)
		allReleaseImageMappings, err := DB.GetAllReleaseImgMappings(db, limit, offset)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = utils.SendAsJSON(w, allReleaseImageMappings)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Accepts a JSON payload of a new image and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 Created
func createReleaseImageMappingHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var receivedReleaseImageMapping objects.ReleaseImageMapping
		err := utils.DecodeJSONObject(w, r, &receivedReleaseImageMapping)
		if err != nil || receivedReleaseImageMapping.ReleaseId == nil || receivedReleaseImageMapping.ImageId == nil {
			log.Error(err)
			return
		}
		err = DB.AddReleaseImgMapping(db, receivedReleaseImageMapping)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		newReleaseImageMapping, err := DB.GetReleaseImageMapping(db, receivedReleaseImageMapping)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"releaseImageMappingId": newReleaseImageMapping.Id,
			"releaseId":             *newReleaseImageMapping.ReleaseId,
			"imageId":               *newReleaseImageMapping.ImageId,
		}).Info("release_image_mapping has been successfully created")
		w.WriteHeader(http.StatusCreated)
		err = utils.SendAsJSON(w, newReleaseImageMapping)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Deletes the releaseImageMapping and responds with an empty payload
//
// Success Code: 204 No Content
func deleteReleaseImageMappingHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var receivedReleaseImageMapping objects.ReleaseImageMapping
		err := utils.DecodeJSONObject(w, r, &receivedReleaseImageMapping)
		if err != nil || receivedReleaseImageMapping.ReleaseId == nil || receivedReleaseImageMapping.ImageId == nil {
			utils.HttpJSONError(w, err.Error(), http.StatusBadRequest)
			log.Error(err)
			return
		}
		err = DB.DeleteReleaseImgMapping(db, receivedReleaseImageMapping)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"releaseImageMappingId": receivedReleaseImageMapping.Id,
			"releaseId":             *receivedReleaseImageMapping.ReleaseId,
			"imageId":               *receivedReleaseImageMapping.ImageId,
		}).Info("release_image_mapping has been successfully deleted")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	return http.HandlerFunc(fn)
}
