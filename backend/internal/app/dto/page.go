package dto

import "github.com/trysourcetool/sourcetool/backend/internal/domain/page"

type ListPagesInput struct {
	EnvironmentID string
}

type Page struct {
	ID             string
	OrganizationID string
	EnvironmentID  string
	APIKeyID       string
	Name           string
	Route          string
	Path           []int32
	CreatedAt      int64
	UpdatedAt      int64
}

func PageFromModel(page *page.Page) *Page {
	if page == nil {
		return nil
	}

	return &Page{
		ID:             page.ID.String(),
		OrganizationID: page.OrganizationID.String(),
		EnvironmentID:  page.EnvironmentID.String(),
		APIKeyID:       page.APIKeyID.String(),
		Name:           page.Name,
		Route:          page.Route,
		Path:           page.Path,
		CreatedAt:      page.CreatedAt.Unix(),
		UpdatedAt:      page.UpdatedAt.Unix(),
	}
}

type ListPagesOutput struct {
	Pages      []*Page
	Groups     []*Group
	GroupPages []*GroupPage
	Users      []*User
	UserGroups []*UserGroup
}
