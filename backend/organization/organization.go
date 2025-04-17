package organization

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Organization struct {
	ID        uuid.UUID `db:"id"`
	Subdomain *string   `db:"subdomain"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
