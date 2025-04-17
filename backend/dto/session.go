package dto

import (
	"github.com/trysourcetool/sourcetool/backend/session"
)

// Session represents session data in DTOs.
type Session struct {
	ID             string
	OrganizationID string
	UserID         string
	APIKeyID       string
	HostInstanceID string
	CreatedAt      int64
	UpdatedAt      int64
}

// SessionFromModel converts from model.Session to dto.Session.
func SessionFromModel(session *session.Session) *Session {
	if session == nil {
		return nil
	}

	return &Session{
		ID:             session.ID.String(),
		OrganizationID: session.OrganizationID.String(),
		UserID:         session.UserID.String(),
		APIKeyID:       session.APIKeyID.String(),
		HostInstanceID: session.HostInstanceID.String(),
		CreatedAt:      session.CreatedAt.Unix(),
		UpdatedAt:      session.UpdatedAt.Unix(),
	}
}
