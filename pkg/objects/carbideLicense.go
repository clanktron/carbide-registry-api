package objects

import (
	"carbide-images-api/pkg/license"
	"crypto/rsa"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
)

type CarbideLicense struct {
	CustomerID     *string
	DaysTillExpiry *int
	NodeCount      *int

	License  *string
	Token    *armcontainerregistry.Token
	Password *armcontainerregistry.TokenPassword
}

func (l CarbideLicense) Valid(pubkeys []*rsa.PublicKey) bool {
	if err := license.ValidateCarbideLicense(l.License, pubkeys); err != nil {
		return false
	}
	return true
}
