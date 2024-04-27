package database

import (
	"carbide-registry-api/pkg/objects"
	"database/sql"
	"errors"
	"fmt"
)

func GetImagebyId(db *sql.DB, imageId int32) (objects.Image, error) {
	var image objects.Image
	err := db.QueryRow(`SELECT * FROM images WHERE id = ?`, imageId).Scan(&image.Id, &image.ImageName, &image.ImageSigned, &image.TrivySigned, &image.TrivyValid, &image.SbomSigned, &image.SbomValid, &image.LastScannedAt, &image.CreatedAt, &image.UpdatedAt)
	if err != nil {
		return image, err
	}
	image.Releases, err = GetAllReleasesforImage(db, imageId, 9999999, 0)
	if err != nil {
		return image, err
	}
	return image, nil
}

func AddImage(db *sql.DB, newImage objects.Image) error {
	const requiredField string = "missing field \"%s\" required when creating a new image"
	// const sqlError string = "error creating new image: %w"
	var (
		imageName     sql.NullString
		imageSigned   sql.NullBool
		trivySigned   sql.NullBool
		trivyValid    sql.NullBool
		sbomSigned    sql.NullBool
		sbomValid     sql.NullBool
		lastScannedAt sql.NullTime
	)

	stmt, err := db.Prepare(`
		INSERT INTO images (image_name, image_signed, trivy_signed, trivy_valid, sbom_signed, sbom_valid, last_scanned_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if newImage.ImageName == nil {
		errMsg := fmt.Sprintf(requiredField, "ImageName")
		return errors.New(errMsg)
	} else {
		imageName.String = *newImage.ImageName
		imageName.Valid = true
	}
	if newImage.ImageSigned != nil {
		imageSigned.Bool = *newImage.ImageSigned
		imageSigned.Valid = true
	}
	if newImage.TrivySigned != nil {
		trivySigned.Bool = *newImage.TrivySigned
		trivySigned.Valid = true
	}
	if newImage.TrivyValid != nil {
		trivyValid.Bool = *newImage.TrivyValid
		trivyValid.Valid = true
	}
	if newImage.SbomSigned != nil {
		sbomSigned.Bool = *newImage.SbomSigned
		sbomSigned.Valid = true
	}
	if newImage.SbomValid != nil {
		sbomValid.Bool = *newImage.SbomValid
		sbomValid.Valid = true
	}
	if newImage.LastScannedAt != nil {
		lastScannedAt.Time = *newImage.LastScannedAt
		lastScannedAt.Valid = true
	}
	_, err = stmt.Exec(
		imageName,
		imageSigned,
		trivySigned,
		trivyValid,
		sbomSigned,
		sbomValid,
		lastScannedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetImagebyName(db *sql.DB, imageName string) (objects.Image, error) {
	var image objects.Image
	err := db.QueryRow(`SELECT * FROM images WHERE image_name = ?`, imageName).Scan(&image.Id, &image.ImageName, &image.ImageSigned, &image.TrivySigned, &image.TrivyValid, &image.SbomSigned, &image.SbomValid, &image.LastScannedAt, &image.CreatedAt, &image.UpdatedAt)
	if err != nil {
		return image, err
	}
	image.Releases, err = GetAllReleasesforImage(db, image.Id, 9999999, 0)
	if err != nil {
		return image, err
	}
	return image, nil
}

func GetAllImages(db *sql.DB, limit int, offset int) ([]objects.Image, error) {
	var images []objects.Image
	rows, err := db.Query(`SELECT * FROM images LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		images = nil
		return images, err
	}
	defer rows.Close()
	for rows.Next() {
		var image objects.Image
		err = rows.Scan(&image.Id, &image.ImageName, &image.ImageSigned, &image.TrivySigned, &image.TrivyValid, &image.SbomSigned, &image.SbomValid, &image.LastScannedAt, &image.CreatedAt, &image.UpdatedAt)
		if err != nil {
			images = nil
			return images, err
		}
		images = append(images, image)
	}
	if err = rows.Err(); err != nil {
		images = nil
		return images, err
	}
	return images, nil
}

func UpdateImage(db *sql.DB, updatedImage objects.Image) error {
	const requiredField string = "missing field \"%s\" required when updating an image"
	// const sqlError string = "error updating image: %w"
	var (
		imageid       sql.NullInt32
		imageSigned   sql.NullBool
		trivySigned   sql.NullBool
		trivyValid    sql.NullBool
		sbomSigned    sql.NullBool
		sbomValid     sql.NullBool
		lastScannedAt sql.NullTime
	)
	stmt, err := db.Prepare(`
		UPDATE images 
		SET image_signed = ?, 
		trivy_signed = ?, 
		trivy_valid = ?, 
		sbom_signed = ?, 
		sbom_valid = ?, 
		last_scanned_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if updatedImage.Id == 0 {
		errMsg := fmt.Sprintf(requiredField, "Id")
		return errors.New(errMsg)
	} else {
		imageid.Int32 = updatedImage.Id
		imageid.Valid = true
	}
	if updatedImage.ImageSigned != nil {
		imageSigned.Bool = *updatedImage.ImageSigned
		imageSigned.Valid = true
	}
	if updatedImage.TrivySigned != nil {
		trivySigned.Bool = *updatedImage.TrivySigned
		trivySigned.Valid = true
	}
	if updatedImage.TrivyValid != nil {
		trivyValid.Bool = *updatedImage.TrivyValid
		trivyValid.Valid = true
	}
	if updatedImage.SbomSigned != nil {
		sbomSigned.Bool = *updatedImage.SbomSigned
		sbomSigned.Valid = true
	}
	if updatedImage.SbomValid != nil {
		sbomValid.Bool = *updatedImage.SbomValid
		sbomValid.Valid = true
	}
	if updatedImage.LastScannedAt != nil {
		lastScannedAt.Time = *updatedImage.LastScannedAt
		lastScannedAt.Valid = true
	}
	_, err = stmt.Exec(
		imageSigned,
		trivySigned,
		trivyValid,
		sbomSigned,
		sbomValid,
		lastScannedAt,
		imageid,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteImage(db *sql.DB, id int32) error {
	_, err := db.Exec(`DELETE FROM images WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return nil
}

func GetImageWithoutReleases(db *sql.DB, imageId int32) (objects.Image, error) {
	var retrievedImage objects.Image
	const sqlError string = "error fetching image: %w"
	err := db.QueryRow(`SELECT * FROM images WHERE id = ?`, imageId).Scan(&retrievedImage.Id, &retrievedImage.ImageName, &retrievedImage.ImageSigned, &retrievedImage.TrivySigned, &retrievedImage.TrivyValid, &retrievedImage.SbomSigned, &retrievedImage.SbomValid, &retrievedImage.LastScannedAt, &retrievedImage.CreatedAt, &retrievedImage.UpdatedAt)
	if err != nil {
		return retrievedImage, fmt.Errorf(sqlError, err)
	}
	return retrievedImage, nil
}

func GetAllImagesforRelease(db *sql.DB, releaseId int32, limit int, offset int) ([]objects.Image, error) {
	var fetchedImages []objects.Image
	var releaseImageMappings []objects.ReleaseImageMapping
	releaseImageMappings, err := GetImgMappings(db, releaseId, limit, offset)
	if err != nil {
		return fetchedImages, err
	}
	for _, releaseImageMapping := range releaseImageMappings {
		image, err := GetImageWithoutReleases(db, *releaseImageMapping.ImageId)
		if err != nil {
			return fetchedImages, err
		}
		fetchedImages = append(fetchedImages, image)
	}
	return fetchedImages, nil
}
