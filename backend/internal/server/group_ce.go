//go:build !ee
// +build !ee

package server

import (
	"errors"
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func (s *Server) getGroup(w http.ResponseWriter, r *http.Request) error {
	return errdefs.ErrPermissionDenied(errors.New("group functionality is not available in CE version"))
}

func (s *Server) listGroups(w http.ResponseWriter, r *http.Request) error {
	return errdefs.ErrPermissionDenied(errors.New("group functionality is not available in CE version"))
}

func (s *Server) createGroup(w http.ResponseWriter, r *http.Request) error {
	return errdefs.ErrPermissionDenied(errors.New("group functionality is not available in CE version"))
}

func (s *Server) updateGroup(w http.ResponseWriter, r *http.Request) error {
	return errdefs.ErrPermissionDenied(errors.New("group functionality is not available in CE version"))
}

func (s *Server) deleteGroup(w http.ResponseWriter, r *http.Request) error {
	return errdefs.ErrPermissionDenied(errors.New("group functionality is not available in CE version"))
}
