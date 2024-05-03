package api

import (
	"carbide-images-api/pkg/api/utils"
	"carbide-images-api/pkg/azure"
	"carbide-images-api/pkg/license"
	"carbide-images-api/pkg/objects"
	"crypto/rsa"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/ebauman/golicense/pkg/certificate"
	log "github.com/sirupsen/logrus"
)

var PRIVATEKEY *rsa.PrivateKey

func init() {
	var GOLICENSE_KEY = os.Getenv("GOLICENSE_KEY")
	var err error
	PRIVATEKEY, err = certificate.PEMToPrivateKey([]byte(GOLICENSE_KEY))
	if err != nil {
		panic(err)
	}
}

func createCarbideAccountHandler(clientFactory *armcontainerregistry.ClientFactory) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var newLicense objects.CarbideLicense
		err := utils.DecodeJSONObject(w, r, &newLicense)
		if err != nil {
			log.Error(err)
			utils.HttpJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if newLicense.CustomerID == nil || newLicense.DaysTillExpiry == nil || newLicense.NodeCount == nil {
			log.Errorf("missing customerID, daysTillExpiry, or nodeCount")
			utils.HttpJSONError(w, "missing customerID, daysTillExpiry, or nodeCount", http.StatusBadRequest)
			return
		}
		expiry := time.Now().Add(time.Hour * 24 * time.Duration(*newLicense.DaysTillExpiry))
		newLicense.License, err = license.CreateCarbideLicense(PRIVATEKEY, *newLicense.NodeCount, *newLicense.CustomerID, expiry)
		if err != nil {
			log.Error(err)
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newLicense.Token, newLicense.Password, err = azure.CreateCarbideAccount(clientFactory, *newLicense.CustomerID, expiry)
		if err != nil {
			log.Error(err)
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		err = utils.SendAsJSON(w, newLicense)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}
