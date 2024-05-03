package license

import (
	"crypto/rand"
	"crypto/rsa"
	mathrand "math/rand"
	"testing"
	"time"
)

const CUSTOMER_ID_PREFIX = "DELETE-ME-"

func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[mathrand.Intn(len(charset))]
	}
	return string(b)
}

func TestCreateCarbideLicense(t *testing.T) {
	nodeCount := 25
	daysTillExpiry := 35
	customerID := CUSTOMER_ID_PREFIX + randString(5)
	expiry := time.Now().Add(time.Hour * 24 * time.Duration(daysTillExpiry))
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	license, err := CreateCarbideLicense(privateKey, nodeCount, customerID, expiry)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated license: %s", *license)
}
