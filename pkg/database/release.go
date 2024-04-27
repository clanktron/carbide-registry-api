package database

import (
	"carbide-registry-api/pkg/objects"
	"database/sql"
	"errors"
	"fmt"
)

func AddRelease(db *sql.DB, newRelease objects.Release) error {
	const requiredField string = "missing field \"%s\" required when creating a new release"
	const sqlError string = "error creating new release: %w"
	if newRelease.ProductId == nil {
		errMsg := fmt.Sprintf(requiredField, "Product Id")
		return errors.New(errMsg)
	}
	if newRelease.Name == nil {
		errMsg := fmt.Sprintf(requiredField, "Name")
		return errors.New(errMsg)
	}
	if newRelease.TarballLink == nil {
		_, err := db.Exec(
			"INSERT INTO releases (product_id, name) VALUES (?, ?)",
			*newRelease.ProductId, *newRelease.Name)
		if err != nil {
			return fmt.Errorf(sqlError, err)
		}
	} else {
		_, err := db.Exec(
			"INSERT INTO releases (product_id, name, tarball_link) VALUES (?, ?, ?)",
			*newRelease.ProductId, *newRelease.Name, *newRelease.TarballLink)
		if err != nil {
			return fmt.Errorf(sqlError, err)
		}
	}
	return nil
}

func GetRelease(db *sql.DB, release objects.Release) (objects.Release, error) {
	const requiredField string = "missing field \"%s\" required when retrieving a release"
	const sqlError string = "error finding release: %w"
	var retrievedRelease objects.Release
	if release.ProductId == nil {
		errMsg := fmt.Sprintf(requiredField, "Product Id")
		return retrievedRelease, errors.New(errMsg)
	}
	if release.Name == nil {
		errMsg := fmt.Sprintf(requiredField, "Name")
		return retrievedRelease, errors.New(errMsg)
	}
	err := db.QueryRow(
		`SELECT * FROM releases WHERE name = ? AND product_id = ?`, *release.Name, *release.ProductId).Scan(
		&retrievedRelease.Id, &retrievedRelease.ProductId, &retrievedRelease.Name, &retrievedRelease.TarballLink, &retrievedRelease.CreatedAt, &retrievedRelease.UpdatedAt)
	if err != nil {
		return retrievedRelease, fmt.Errorf(sqlError, err)
	}
	retrievedRelease.Images, err = GetAllImagesforRelease(db, retrievedRelease.Id, 9999999, 0)
	if err != nil {
		return retrievedRelease, err
	}
	return retrievedRelease, nil
}

func GetAllReleasesforProduct(db *sql.DB, productName string, limit int, offset int) ([]objects.Release, error) {
	var releases []objects.Release
	product, err := GetProduct(db, productName)
	if err != nil {
		releases = nil
		return releases, err
	}
	rows, err := db.Query(`SELECT * FROM releases WHERE product_id = ? LIMIT ? OFFSET ?`, product.Id, limit, offset)
	if err != nil {
		releases = nil
		return releases, err
	}
	defer rows.Close()
	for rows.Next() {
		var release objects.Release
		err = rows.Scan(&release.Id, &release.ProductId, &release.Name, &release.TarballLink, &release.CreatedAt, &release.UpdatedAt)
		if err != nil {
			releases = nil
			return releases, err
		}
		releases = append(releases, release)
	}
	if err = rows.Err(); err != nil {
		releases = nil
		return releases, err
	}
	return releases, nil
}

func GetAllReleases(db *sql.DB, limit int, offset int) ([]objects.Release, error) {
	var releases []objects.Release
	rows, err := db.Query(`SELECT * FROM releases LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		releases = nil
		return releases, err
	}
	defer rows.Close()
	for rows.Next() {
		var release objects.Release
		err = rows.Scan(&release.Id, &release.ProductId, &release.Name, &release.TarballLink, &release.CreatedAt, &release.UpdatedAt)
		if err != nil {
			releases = nil
			return releases, err
		}
		releases = append(releases, release)
	}
	if err = rows.Err(); err != nil {
		releases = nil
		return releases, err
	}
	return releases, nil
}

func UpdateRelease(db *sql.DB, updatedRelease objects.Release) error {
	const missingField string = "missing field %s (needed to locate release in DB)"
	const sqlError string = "error updating new release: %w"
	if updatedRelease.ProductId == nil {
		errMsg := fmt.Sprintf(missingField, "Product Id")
		return errors.New(errMsg)
	}
	if updatedRelease.Name == nil {
		errMsg := fmt.Sprintf(missingField, "Name")
		return errors.New(errMsg)
	}
	if updatedRelease.TarballLink == nil {
		return errors.New("no new data to update release with")
	} else {
		_, err := db.Exec(
			`UPDATE releases SET tarball_link = ? WHERE name = ? AND product_id = ?`,
			*updatedRelease.TarballLink, *updatedRelease.Name, *updatedRelease.ProductId)
		if err != nil {
			return fmt.Errorf(sqlError, err)
		}
	}
	return nil
}

func DeleteRelease(db *sql.DB, releaseToDelete objects.Release) error {
	const missingField string = "missing field %s (needed to locate release in DB)"
	// const sqlError string = "error updating new release: %w"
	if releaseToDelete.ProductId == nil {
		errMsg := fmt.Sprintf(missingField, "Product Id")
		return errors.New(errMsg)
	}
	if releaseToDelete.Name == nil {
		errMsg := fmt.Sprintf(missingField, "Name")
		return errors.New(errMsg)
	}
	_, err := db.Exec(
		`DELETE FROM releases WHERE name = ? AND product_id = ?`,
		*releaseToDelete.Name, *releaseToDelete.ProductId)
	if err != nil {
		return err
	}
	return nil
}

func GetReleaseWithoutImages(db *sql.DB, release_id int32) (objects.Release, error) {
	var retrievedRelease objects.Release
	const sqlError string = "error fetching release: %w"
	err := db.QueryRow(
		`SELECT * FROM releases WHERE id = ?`, release_id).Scan(
		&retrievedRelease.Id, &retrievedRelease.ProductId, &retrievedRelease.Name, &retrievedRelease.TarballLink, &retrievedRelease.CreatedAt, &retrievedRelease.UpdatedAt)
	if err != nil {
		return retrievedRelease, fmt.Errorf(sqlError, err)
	}
	return retrievedRelease, nil
}

func GetAllReleasesforImage(db *sql.DB, imageId int32, limit int, offset int) ([]objects.Release, error) {
	var fetchedReleases []objects.Release
	var releaseImageMappings []objects.ReleaseImageMapping
	releaseImageMappings, err := GetReleaseMappings(db, imageId, limit, offset)
	if err != nil {
		return fetchedReleases, err
	}
	for _, releaseImageMapping := range releaseImageMappings {
		release, err := GetReleaseWithoutImages(db, *releaseImageMapping.ReleaseId)
		if err != nil {
			return fetchedReleases, err
		}
		fetchedReleases = append(fetchedReleases, release)
	}
	return fetchedReleases, nil
}
