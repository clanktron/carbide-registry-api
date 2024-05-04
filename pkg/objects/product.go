package objects

import (
	"time"
)

type Product struct {
	Id        int32
	Name      *string
	LogoUrl   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
