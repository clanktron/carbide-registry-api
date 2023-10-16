package api

import (
	DB "carbide-api/pkg/database"
	"carbide-api/pkg/objects"
	"database/sql"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Responds with a JSON array of all products in the database
//
// Success Code: 200 OK
func productGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	products, err := DB.GetAllProducts(db)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	products_json, err := json.Marshal(products)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(products_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Accepts a JSON payload of a new product and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 OK
func productPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var created_product objects.Product
	err := json.NewDecoder(r.Body).Decode(&created_product)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = DB.AddProduct(db, created_product)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	created_product, err = DB.GetProduct(db, *created_product.Name)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"product": *created_product.Name,
	}).Info("Product has been successfully created")
	created_product_json, err := json.Marshal(created_product)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(created_product_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Responds with the JSON representation of a product
//
// Success Code: 200 OK
func productGetByName(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	var retrieved_product objects.Product
	retrieved_product, err := DB.GetProduct(db, product_name)
	if err != nil {
		log.Error(err)
		return
	}
	retrieved_product_json, err := json.Marshal(retrieved_product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(retrieved_product_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Responds with the JSON representation of a product
//
// Success Code: 200 OK
func productPutByName(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	var updated_product objects.Product
	err := json.NewDecoder(r.Body).Decode(&updated_product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = DB.UpdateProduct(db, *updated_product.Name, product_name)
	if err != nil {
		log.Error(err)
		return
	}
	updated_product, err = DB.GetProduct(db, *updated_product.Name)
	if err != nil {
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"product": *updated_product.Name,
	}).Info("Product has been successfully updated")
	updated_product_json, err := json.Marshal(updated_product)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(updated_product_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Deletes the product and responds with an empty payload
//
// Success Code: 204 No Content
func productDeleteByName(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	err := DB.DeleteProduct(db, product_name)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"product": product_name,
	}).Info("Product has been successfully deleted")
	w.WriteHeader(http.StatusNoContent)
	return
}
