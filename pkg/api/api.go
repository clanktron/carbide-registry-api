package api

import (
	"carbide-images-api/pkg/api/middleware"
	"database/sql"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
)

func NewRouter(db *sql.DB, clientFactory *armcontainerregistry.ClientFactory) http.Handler {
	mux := http.NewServeMux()
	withAuth := middleware.SessionAuth
	mux.Handle("GET /healthcheck", middleware.Healthcheck())
	mux.Handle("POST /carbide/license", withAuth(createCarbideAccountHandler(clientFactory)))
	mux.Handle("GET /auth", authCheckHandler())
	mux.Handle("POST /auth", loginHandler())
	mux.Handle("GET /product", withAuth(getAllProductsHandler(db)))
	mux.Handle("POST /product", withAuth(createProductHandler(db)))
	mux.Handle("GET /product/{productName}", withAuth(getProductHandler(db)))
	mux.Handle("PUT /product/{productName}", withAuth(updateProductHandler(db)))
	mux.Handle("DELETE /product/{productName}", withAuth(deleteProductHandler(db)))
	mux.Handle("GET /product/{productName}/release", withAuth(getAllReleasesHandler(db)))
	mux.Handle("POST /product/{productName}/release", withAuth(createReleaseHandler(db)))
	mux.Handle("GET /product/{productName}/release/{releaseName}", withAuth(getReleaseHandler(db)))
	mux.Handle("PUT /product/{productName}/release/{releaseName}", withAuth(updateReleaseHandler(db)))
	mux.Handle("DELETE /product/{proudctName}/release/{releaseName}", withAuth(deleteReleaseHandler(db)))
	mux.Handle("GET /image", withAuth(getAllImagesHandler(db)))
	mux.Handle("POST /image", withAuth(createImageHandler(db)))
	mux.Handle("GET /image/{imageID}", withAuth(getImageHandler(db)))
	mux.Handle("PUT /image/{imageID}", withAuth(updateImageHandler(db)))
	mux.Handle("DELETE /image/{imageID}", withAuth(deleteImageHandler(db)))
	mux.Handle("GET /releaseImageMapping", withAuth(getAllReleaseImageMappingsHandler(db)))
	mux.Handle("POST /releaseImageMapping", withAuth(createReleaseImageMappingHandler(db)))
	mux.Handle("DELETE /releaseImageMapping", withAuth(deleteReleaseImageMappingHandler(db)))
	withCors := middleware.CORS
	newMux := withCors(mux)
	return newMux
}
