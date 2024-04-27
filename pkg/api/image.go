package api

import (
	"carbide-registry-api/pkg/api/utils"
	DB "carbide-registry-api/pkg/database"
	"carbide-registry-api/pkg/objects"
	"database/sql"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// returns imageID if a valid int32, or -1 if not
func imageIDFromPath(w http.ResponseWriter, r *http.Request) int32 {
	imageID64, err := strconv.ParseInt(r.PathValue("imageID"), 10, 32)
	if err != nil {
		utils.HttpJSONError(w, "invalid image ID", http.StatusBadRequest)
		log.Error(err)
		return -1
	}
	if imageID64 < -2147483648 || imageID64 > 2147483647 {
		utils.HttpJSONError(w, "invalid image ID", http.StatusBadRequest)
		log.Error(err)
		return -1
	}
	imageID := int32(imageID64)
	return imageID
}

// Responds with a JSON array of all images in the database
//
// Success Code: 200 OK
func getAllImagesHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		limit, offset := utils.GetLimitAndOffset(r)
		images, err := DB.GetAllImages(db, limit, offset)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = utils.SendAsJSON(w, images)
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
func createImageHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var newImage objects.Image
		err := utils.DecodeJSONObject(w, r, &newImage)
		if err != nil {
			log.Error(err)
		}
		if newImage.ImageName == nil {
			utils.HttpJSONError(w, "missing image name", http.StatusBadRequest)
			log.Error(err)
			return
		}
		err = DB.AddImage(db, newImage)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		createdImage, err := DB.GetImagebyName(db, *newImage.ImageName)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"image": *createdImage.ImageName,
		}).Info("image has been successfully created")
		w.WriteHeader(http.StatusCreated)
		err = utils.SendAsJSON(w, createdImage)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Responds with the JSON representation of an image (includes associated releases)
//
// Success Code: 200 OK
func getImageHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		imageID := imageIDFromPath(w, r)
		if imageID == -1 {
			return
		}
		image, err := DB.GetImagebyId(db, imageID)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = utils.SendAsJSON(w, image)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Accepts a JSON payload of the updated image and responds with the new JSON object after it's been successfully updated in the database
//
// Success Code: 200 OK
func updateImageHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var updatedImage objects.Image
		err := utils.DecodeJSONObject(w, r, &updatedImage)
		if err != nil {
			log.Error(err)
			return
		}
		// image id cannot be overwritten with json payload
		imageID := imageIDFromPath(w, r)
		if imageID == -1 {
			return
		}
		updatedImage.Id = imageID
		err = DB.UpdateImage(db, updatedImage)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		updatedImage, err = DB.GetImagebyId(db, updatedImage.Id)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"image": *updatedImage.ImageName,
		}).Info("image has been successfully updated")
		w.WriteHeader(http.StatusOK)
		err = utils.SendAsJSON(w, updatedImage)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Deletes the image and responds with an empty payload
//
// Success Code: 204 No Content
func deleteImageHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		imageID := imageIDFromPath(w, r)
		if imageID == -1 {
			return
		}
		err := DB.DeleteImage(db, imageID)
		if err != nil {
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"image": imageID,
		}).Info("image has been successfully deleted")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	return http.HandlerFunc(fn)
}
