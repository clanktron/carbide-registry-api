package objects

import (
	"errors"
	"fmt"
	"time"
)

type Product struct {
	Id        int32
	Name      *string
	LogoUrl   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Product) Validate() error {
	const requiredField string = "missing field \"%s\" required for product"
	if p.Name == nil {
		errMsg := fmt.Sprintf(requiredField, "Name")
		return errors.New(errMsg)
	}
	if p.LogoUrl == nil {
		errMsg := fmt.Sprintf(requiredField, "LogoUrl")
		return errors.New(errMsg)
	}
	return nil
}
