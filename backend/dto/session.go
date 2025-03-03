package dto

import (
	"github.com/gofrs/uuid/v5"
	"github.com/trysourcetool/sourcetool/backend/model"
)

// Session represents session data in DTOs
type Session struct {
	ID             string
	OrganizationID string
	UserID         string
	PageID         string
	HostInstanceID string
	CreatedAt      int64
	UpdatedAt      int64
}

// SessionFromModel converts from model.Session to dto.Session
func SessionFromModel(session *model.Session) *Session {
	if session == nil {
		return nil
	}

	return &Session{
		ID:             session.ID.String(),
		OrganizationID: session.OrganizationID.String(),
		UserID:         session.UserID.String(),
		PageID:         session.PageID.String(),
		HostInstanceID: session.HostInstanceID.String(),
		CreatedAt:      session.CreatedAt.Unix(),
		UpdatedAt:      session.UpdatedAt.Unix(),
	}
}

// ToSessionID converts string ID to uuid.UUID
func ToSessionID(id string) (uuid.UUID, error) {
	return uuid.FromString(id)
}
