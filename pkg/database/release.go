package database

import (
	"carbide-api/pkg/objects"
	"database/sql"
	"errors"
	"fmt"
)

func AddRelease(db *sql.DB, new_release objects.Release) error {
	const required_field string = "Missing field \"%s\" required when creating a new release"
	const sql_error string = "Error creating new release: %w"
	if new_release.ProductId == nil {
		err_msg := fmt.Sprintf(required_field, "Product Id")
		return errors.New(err_msg)
	}
	if new_release.Name == nil {
		err_msg := fmt.Sprintf(required_field, "Name")
		return errors.New(err_msg)
	}
	if new_release.TarballLink == nil {
		_, err := db.Exec(
			"INSERT INTO releases (product_id, name) VALUES (?, ?)",
			*new_release.ProductId, *new_release.Name)
		if err != nil {
			return fmt.Errorf(sql_error, err)
		}
	} else {
		_, err := db.Exec(
			"INSERT INTO releases (product_id, name, tarball_link) VALUES (?, ?, ?)",
			*new_release.ProductId, *new_release.Name, *new_release.TarballLink)
		if err != nil {
			return fmt.Errorf(sql_error, err)
		}
	}
	return nil
}

func GetRelease(db *sql.DB, release objects.Release) (objects.Release, error) {
	const required_field string = "Missing field \"%s\" required when retrieving a release"
	const sql_error string = "Error finding release: %w"
	var retrieved_release objects.Release
	if release.ProductId == nil {
		err_msg := fmt.Sprintf(required_field, "Product Id")
		return retrieved_release, errors.New(err_msg)
	}
	if release.Name == nil {
		err_msg := fmt.Sprintf(required_field, "Name")
		return retrieved_release, errors.New(err_msg)
	}
	err := db.QueryRow(
		`SELECT * FROM releases WHERE name = ? AND product_id = ?`, *release.Name, *release.ProductId).Scan(
		&retrieved_release.Id, &retrieved_release.ProductId, &retrieved_release.Name, &retrieved_release.TarballLink, &retrieved_release.CreatedAt, &retrieved_release.UpdatedAt)
	if err != nil {
		return retrieved_release, fmt.Errorf(sql_error, err)
	}
	retrieved_release.Images, err = GetAllImagesforRelease(db, retrieved_release.Id)
	if err != nil {
		return retrieved_release, err
	}
	return retrieved_release, nil
}

func GetAllReleasesforProduct(db *sql.DB, product_name string) ([]objects.Release, error) {

	product, err := GetProduct(db, product_name)
	product_id := product.Id

	var releases []objects.Release
	rows, err := db.Query(`SELECT * FROM releases WHERE product_id = ?`, product_id)
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

func GetAllReleases(db *sql.DB) ([]objects.Release, error) {
	var releases []objects.Release
	rows, err := db.Query(`SELECT * FROM releases`)
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

func UpdateRelease(db *sql.DB, updated_release objects.Release) error {
	const missing_field string = "Missing field %s (needed to locate release in DB)"
	const sql_error string = "Error updating new release: %w"
	if updated_release.ProductId == nil {
		err_msg := fmt.Sprintf(missing_field, "Product Id")
		return errors.New(err_msg)
	}
	if updated_release.Name == nil {
		err_msg := fmt.Sprintf(missing_field, "Name")
		return errors.New(err_msg)
	}
	if updated_release.TarballLink == nil {
		return errors.New("No new data to update release with")
	} else {
		_, err := db.Exec(
			`UPDATE releases SET tarball_link = ? WHERE name = ? AND product_id = ?`,
			*updated_release.TarballLink, *updated_release.Name, *updated_release.ProductId)
		if err != nil {
			return fmt.Errorf(sql_error, err)
		}
	}
	return nil
}

func DeleteRelease(db *sql.DB, release_to_delete objects.Release) error {
	const missing_field string = "Missing field %s (needed to locate release in DB)"
	const sql_error string = "Error updating new release: %w"
	if release_to_delete.ProductId == nil {
		err_msg := fmt.Sprintf(missing_field, "Product Id")
		return errors.New(err_msg)
	}
	if release_to_delete.Name == nil {
		err_msg := fmt.Sprintf(missing_field, "Name")
		return errors.New(err_msg)
	}
	_, err := db.Exec(
		`DELETE FROM releases WHERE name = ? AND product_id = ?`,
		*release_to_delete.Name, *release_to_delete.ProductId)
	if err != nil {
		return err
	}
	return nil
}

func GetReleaseWithoutImages(db *sql.DB, release_id int32) (objects.Release, error) {
	var retrieved_release objects.Release
	const sql_error string = "Error fetching release: %w"
	err := db.QueryRow(
		`SELECT * FROM releases WHERE id = ?`, release_id).Scan(
		&retrieved_release.Id, &retrieved_release.ProductId, &retrieved_release.Name, &retrieved_release.TarballLink, &retrieved_release.CreatedAt, &retrieved_release.UpdatedAt)
	if err != nil {
		return retrieved_release, fmt.Errorf(sql_error, err)
	}
	return retrieved_release, nil

}

func GetAllReleasesforImage(db *sql.DB, image_id int32) ([]objects.Release, error) {

	var fetched_releases []objects.Release

	var release_img_mappings []objects.Release_Image_Mapping
	release_img_mappings, err := GetReleaseMappings(db, image_id)
	if err != nil {
		return fetched_releases, err
	}

	for _, release_image_mapping := range release_img_mappings {
		release, err := GetReleaseWithoutImages(db, *release_image_mapping.ReleaseId)
		if err != nil {
			return fetched_releases, err
		}
		fetched_releases = append(fetched_releases, release)
	}

	return fetched_releases, nil
}
