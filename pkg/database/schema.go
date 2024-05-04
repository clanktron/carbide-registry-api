package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Init(db_user string, db_pass string, db_host string, db_port string, db_name string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db_user, db_pass, db_host, db_port, db_name)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return db, err
	} else if err = db.Ping(); err != nil {
		return db, err
	}
	return db, err
}

func SchemaInit(db *sql.DB) error {
	err := ensureProductTable(db)
	if err != nil {
		return err
	}
	err = ensureReleasesTable(db)
	if err != nil {
		return err
	}
	err = ensureImagesTable(db)
	if err != nil {
		return err
	}
	err = ensureUsersTable(db)
	if err != nil {
		return err
	}
	err = ensureReleaseImageMappingTable(db)
	if err != nil {
		return err
	}
	return nil
}

func ensureUsersTable(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		  username VARCHAR(255) NOT NULL UNIQUE,
		  password VARCHAR(255) NOT NULL UNIQUE,
		  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`); err != nil {
		return err
	}
	return nil
}

func ensureProductTable(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS product (
		  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		  name VARCHAR(255) NOT NULL,
		  logo_url VARCHAR(255) NOT NULL,
		  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`); err != nil {
		return err
	}
	return nil
}

func ensureReleasesTable(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS releases (
		  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		  product_id INT NOT NULL,
		  name VARCHAR(255) NOT NULL,
		  tarball_link VARCHAR(255),
		  CONSTRAINT fk_release_product FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE,
		  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`); err != nil {
		return err
	}
	return nil
}

func ensureImagesTable(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS images (
		  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		  image_name VARCHAR(255) NOT NULL,
		  image_signed BOOLEAN DEFAULT false,
		  trivy_signed BOOLEAN DEFAULT false,
		  trivy_valid BOOLEAN DEFAULT false,
		  sbom_signed BOOLEAN DEFAULT false,
		  sbom_valid BOOLEAN DEFAULT false,
		  last_scanned_at DATETIME DEFAULT NULL,
		  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`); err != nil {
		return err
	}
	return nil
}

func ensureReleaseImageMappingTable(db *sql.DB) error {
	if _, err := db.Exec(`
   	CREATE TABLE IF NOT EXISTS release_image_mapping (
   	  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   	  release_id INT NOT NULL,
   	  image_id INT NOT NULL,
   	  CONSTRAINT fk_releases_map FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE,
   	  CONSTRAINT fk_images_map FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
   	  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   	  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
   	)`); err != nil {
		return err
	}
	return nil
}
