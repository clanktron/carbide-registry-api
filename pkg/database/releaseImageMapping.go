package database

import (
	"carbide-registry-api/pkg/objects"
	"database/sql"
	"errors"
	"fmt"
)

func GetReleaseImageMappingbyId(db *sql.DB, releaseImageMappingId int32) (objects.ReleaseImageMapping, error) {
	var releaseImageMapping objects.ReleaseImageMapping
	err := db.QueryRow(`SELECT * FROM release_image_mapping WHERE id = ?`, releaseImageMappingId).Scan(
		&releaseImageMapping.Id,
		&releaseImageMapping.ReleaseId,
		&releaseImageMapping.ImageId,
		&releaseImageMapping.CreatedAt,
		&releaseImageMapping.UpdatedAt,
	)
	if err != nil {
		return releaseImageMapping, err
	}
	return releaseImageMapping, nil
}

func GetReleaseImageMapping(db *sql.DB, releaseImageMapping objects.ReleaseImageMapping) (objects.ReleaseImageMapping, error) {
	const requiredField string = "Missing field \"%s\" required when retrieving a releaseImageMapping"
	// const sqlError string = "Error retrieving releaseImageMapping: %w"
	var retrievedReleaseImageMapping objects.ReleaseImageMapping
	if releaseImageMapping.ReleaseId == nil {
		errMsg := fmt.Sprintf(requiredField, "ReleaseId")
		return retrievedReleaseImageMapping, errors.New(errMsg)
	}
	if releaseImageMapping.ImageId == nil {
		errMsg := fmt.Sprintf(requiredField, "ImageId")
		return retrievedReleaseImageMapping, errors.New(errMsg)
	}
	err := db.QueryRow(`SELECT * FROM release_image_mapping WHERE release_id = ? AND image_id = ?`, *releaseImageMapping.ReleaseId, *releaseImageMapping.ImageId).Scan(
		&retrievedReleaseImageMapping.Id,
		&retrievedReleaseImageMapping.ReleaseId,
		&retrievedReleaseImageMapping.ImageId,
		&retrievedReleaseImageMapping.CreatedAt,
		&retrievedReleaseImageMapping.UpdatedAt,
	)
	if err != nil {
		return retrievedReleaseImageMapping, err
	}
	return retrievedReleaseImageMapping, nil
}

func AddReleaseImgMapping(db *sql.DB, newReleaseImgMapping objects.ReleaseImageMapping) error {
	const requiredField string = "Missing field \"%s\" required when creating a releaseImageMapping"
	// const sqlError string = "Error creating releaseImageMapping: %w"
	if newReleaseImgMapping.ReleaseId == nil {
		errMsg := fmt.Sprintf(requiredField, "ReleaseId")
		return errors.New(errMsg)
	}
	if newReleaseImgMapping.ImageId == nil {
		errMsg := fmt.Sprintf(requiredField, "ImageId")
		return errors.New(errMsg)
	}
	_, err := db.Exec("INSERT INTO release_image_mapping (release_id, image_id) VALUES (?, ?)",
		*newReleaseImgMapping.ReleaseId, *newReleaseImgMapping.ImageId)
	if err != nil {
		return err
	}
	return nil
}

func GetImgMappings(db *sql.DB, releaseId int32, limit int, offset int) ([]objects.ReleaseImageMapping, error) {
	var releaseImgMappings []objects.ReleaseImageMapping
	rows, err := db.Query(`SELECT * FROM release_image_mapping WHERE release_id = ? LIMIT ? OFFSET ?`, releaseId, limit, offset)
	if err != nil {
		releaseImgMappings = nil
		return releaseImgMappings, err
	}
	defer rows.Close()
	for rows.Next() {
		var releaseImgMapping objects.ReleaseImageMapping
		err = rows.Scan(&releaseImgMapping.Id, &releaseImgMapping.ReleaseId, &releaseImgMapping.ImageId, &releaseImgMapping.CreatedAt, &releaseImgMapping.UpdatedAt)
		if err != nil {
			releaseImgMappings = nil
			return releaseImgMappings, err
		}
		releaseImgMappings = append(releaseImgMappings, releaseImgMapping)
	}
	if err = rows.Err(); err != nil {
		releaseImgMappings = nil
		return releaseImgMappings, err
	}
	return releaseImgMappings, nil
}

func GetReleaseMappings(db *sql.DB, imageId int32, limit int, offset int) ([]objects.ReleaseImageMapping, error) {
	var releaseImgMappings []objects.ReleaseImageMapping
	rows, err := db.Query(`SELECT * FROM release_image_mapping WHERE image_id = ? LIMIT ? OFFSET ?`, imageId, limit, offset)
	if err != nil {
		releaseImgMappings = nil
		return releaseImgMappings, err
	}
	defer rows.Close()
	for rows.Next() {
		var releaseImgMapping objects.ReleaseImageMapping
		err = rows.Scan(&releaseImgMapping.Id, &releaseImgMapping.ReleaseId, &releaseImgMapping.ImageId, &releaseImgMapping.CreatedAt, &releaseImgMapping.UpdatedAt)
		if err != nil {
			releaseImgMappings = nil
			return releaseImgMappings, err
		}
		releaseImgMappings = append(releaseImgMappings, releaseImgMapping)
	}
	if err = rows.Err(); err != nil {
		releaseImgMappings = nil
		return releaseImgMappings, err
	}
	return releaseImgMappings, nil
}

func GetAllReleaseImgMappings(db *sql.DB, limit int, offset int) ([]objects.ReleaseImageMapping, error) {
	var releaseImgMappings []objects.ReleaseImageMapping
	rows, err := db.Query(`SELECT * FROM release_image_mapping LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		releaseImgMappings = nil
		return releaseImgMappings, err
	}
	defer rows.Close()

	for rows.Next() {
		var releaseImgMapping objects.ReleaseImageMapping
		err = rows.Scan(&releaseImgMapping.Id, &releaseImgMapping.ReleaseId, &releaseImgMapping.ImageId, &releaseImgMapping.CreatedAt, &releaseImgMapping.UpdatedAt)
		if err != nil {
			releaseImgMappings = nil
			return releaseImgMappings, err
		}
		releaseImgMappings = append(releaseImgMappings, releaseImgMapping)
	}
	if err = rows.Err(); err != nil {
		releaseImgMappings = nil
		return releaseImgMappings, err
	}

	return releaseImgMappings, nil
}

func DeleteReleaseImgMappingbyId(db *sql.DB, releaseImgMappingId int32) error {
	_, err := db.Exec(`DELETE FROM release_image_mapping WHERE id = ?`, releaseImgMappingId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteReleaseImgMapping(db *sql.DB, releaseImageMappingToDelete objects.ReleaseImageMapping) error {
	const requiredField string = "Missing field \"%s\" required when deleting a releaseImageMapping"
	// const sqlError string = "Error deleting releaseImageMapping: %w"
	if releaseImageMappingToDelete.ReleaseId == nil {
		errMsg := fmt.Sprintf(requiredField, "ReleaseId")
		return errors.New(errMsg)
	}
	if releaseImageMappingToDelete.ImageId == nil {
		errMsg := fmt.Sprintf(requiredField, "ImageId")
		return errors.New(errMsg)
	}
	_, err := db.Exec(`DELETE FROM release_image_mapping WHERE release_id = ? AND image_id = ?`, releaseImageMappingToDelete.ReleaseId, releaseImageMappingToDelete.ImageId)
	if err != nil {
		return err
	}
	return nil
}
