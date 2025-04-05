package sourcetool

import (
	"strings"

	"github.com/gofrs/uuid/v5"
)

type Router interface {
	Page(relativePath, name string, handler func(UIBuilder) error)
	AccessGroups(groups ...string) Router
	Group(relativePath string) Router
}

type router struct {
	parent       *router
	sourcetool   *Sourcetool
	basePath     string
	namespaceDNS string
	groups       []string
}

func newRouter(st *Sourcetool, namespaceDNS string) Router {
	return &router{
		groups:       make([]string, 0),
		sourcetool:   st,
		namespaceDNS: namespaceDNS,
	}
}

func (r *router) generatePageID(fullPath string) uuid.UUID {
	ns := uuid.NewV5(uuid.NamespaceDNS, r.namespaceDNS)
	return uuid.NewV5(ns, fullPath)
}

func (r *router) joinPath(relativePath string) string {
	if !strings.HasPrefix(relativePath, "/") {
		relativePath = "/" + relativePath
	}
	if r.basePath == "" {
		if relativePath == "/" {
			return relativePath
		}
		return strings.TrimSuffix(relativePath, "/")
	}
	basePath := strings.TrimSuffix(r.basePath, "/")
	cleanPath := strings.TrimPrefix(relativePath, "/")
	result := basePath + "/" + cleanPath
	if result == "/" {
		return result
	}
	return strings.TrimSuffix(result, "/")
}

func removeDuplicates(groups []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(groups))

	for _, group := range groups {
		if _, exists := seen[group]; !exists {
			seen[group] = struct{}{}
			result = append(result, group)
		}
	}
	return result
}

func (r *router) collectGroups() []string {
	groups := make([]string, 0)
	current := r
	for current != nil {
		groups = append(groups, current.groups...)
		current = current.parent
	}
	return groups
}

func (r *router) Page(relativePath, name string, handler func(UIBuilder) error) {
	// Skip page creation only for top-level root path
	if relativePath == "/" && r.basePath == "" {
		return
	}

	var fullPath string
	if relativePath == "" {
		if r.basePath == "" {
			fullPath = "/"
		} else {
			fullPath = strings.TrimSuffix(r.basePath, "/")
		}
	} else {
		fullPath = r.joinPath(relativePath)
	}
	pageID := r.generatePageID(fullPath)

	page := &page{
		id:           pageID,
		name:         name,
		route:        fullPath,
		path:         []int{len(r.sourcetool.pages)},
		handler:      handler,
		accessGroups: removeDuplicates(r.collectGroups()),
	}

	r.sourcetool.addPage(pageID, page)
}

func (r *router) AccessGroups(groups ...string) Router {
	if len(groups) > 0 {
		r.groups = append(r.groups, groups...)
	}
	return r
}

func (r *router) Group(relativePath string) Router {
	newRouter := &router{
		parent:       r,
		sourcetool:   r.sourcetool,
		basePath:     r.joinPath(relativePath),
		namespaceDNS: r.namespaceDNS,
		groups:       make([]string, 0),
	}

	return newRouter
}
