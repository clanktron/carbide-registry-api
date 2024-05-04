package objects

import "testing"

func TestProductValidate(t *testing.T) {
	productName := "rancher"
	productLogoUrl := "https://rancher.com"
	validProduct := Product{
		Name:    &productName,
		LogoUrl: &productLogoUrl,
	}
	if err := validProduct.Validate(); err != nil {
		t.Errorf("Valid product stated to be invalid")
	}
	invalidProduct := Product{
		Name:    nil,
		LogoUrl: nil,
	}
	if err := invalidProduct.Validate(); err == nil {
		t.Errorf("Invalid product stated to be valid")
	}
}
