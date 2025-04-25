//go:build !ee
// +build !ee

package server

import (
	"errors"
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func (s *Server) handleGetGroup(w http.ResponseWriter, r *http.Request) error {
	return errdefs.ErrPermissionDenied(errors.New("group functionality is not available in CE version"))
}

func (s *Server) handleListGroups(w http.ResponseWriter, r *http.Request) error {
	return errdefs.ErrPermissionDenied(errors.New("group functionality is not available in CE version"))
}

func (s *Server) handleCreateGroup(w http.ResponseWriter, r *http.Request) error {
	return errdefs.ErrPermissionDenied(errors.New("group functionality is not available in CE version"))
}

func (s *Server) handleUpdateGroup(w http.ResponseWriter, r *http.Request) error {
	return errdefs.ErrPermissionDenied(errors.New("group functionality is not available in CE version"))
}

func (s *Server) handleDeleteGroup(w http.ResponseWriter, r *http.Request) error {
	return errdefs.ErrPermissionDenied(errors.New("group functionality is not available in CE version"))
}
