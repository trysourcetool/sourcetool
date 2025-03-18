package sourcetool

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid/v5"
)

func TestNew(t *testing.T) {
	apiKey := "test_api_key"
	host := "ws://test.trysourcetool.com"
	config := &Config{
		APIKey: apiKey,
		Host:   host,
	}
	st := New(config)

	if st == nil {
		t.Fatal("New returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"APIKey", st.apiKey, apiKey},
		{"Endpoint", st.endpoint, "ws://test.trysourcetool.com/ws"},
		{"Pages length", len(st.pages), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestPage(t *testing.T) {
	pageHandler := func(ui UIBuilder) error { return nil }

	t.Run("Public page", func(t *testing.T) {
		config := &Config{
			APIKey: "test_api_key",
			Host:   "ws://test.trysourcetool.com",
		}
		st := New(config)
		st.Page("/public", "Public Page", pageHandler)

		page := findPageByPath(st.pages, "/public")
		if page == nil {
			t.Fatal("Page not found")
		}

		if len(page.accessGroups) != 0 {
			t.Errorf("Expected no access groups, got %v", page.accessGroups)
		}

		if !page.hasAccess([]string{}) {
			t.Error("Public page should be accessible without groups")
		}
	})

	t.Run("Page with direct access groups", func(t *testing.T) {
		config := &Config{
			APIKey: "test_api_key",
			Host:   "ws://test.trysourcetool.com",
		}
		st := New(config)
		st.AccessGroups("admin")
		st.Page("/admin", "Admin Page", pageHandler)

		page := findPageByPath(st.pages, "/admin")
		if page == nil {
			t.Fatal("Page not found")
		}

		if len(page.accessGroups) != 1 || page.accessGroups[0] != "admin" {
			t.Errorf("Expected [admin] access group, got %v", page.accessGroups)
		}
	})

	t.Run("Group with access groups", func(t *testing.T) {
		config := &Config{
			APIKey: "test_api_key",
			Host:   "ws://test.trysourcetool.com",
		}
		st := New(config)
		api := st.Group("/api")
		api.AccessGroups("api_user")
		api.Page("/users", "Users API", pageHandler)
		api.Page("/posts", "Posts API", pageHandler)

		usersPage := findPageByPath(st.pages, "/api/users")
		postsPage := findPageByPath(st.pages, "/api/posts")

		if usersPage == nil || postsPage == nil {
			t.Fatal("Pages not found")
		}

		if len(usersPage.accessGroups) != 1 || usersPage.accessGroups[0] != "api_user" {
			t.Errorf("Expected [api_user] access group for users page, got %v", usersPage.accessGroups)
		}

		if len(postsPage.accessGroups) != 1 || postsPage.accessGroups[0] != "api_user" {
			t.Errorf("Expected [api_user] access group for posts page, got %v", postsPage.accessGroups)
		}
	})

	t.Run("Nested groups with access groups", func(t *testing.T) {
		config := &Config{
			APIKey: "test_api_key",
			Host:   "ws://test.trysourcetool.com",
		}
		st := New(config)
		users := st.Group("/users")
		users.AccessGroups("admin")
		users.Page("/list", "List users page", pageHandler)
		users.AccessGroups("customer_support").Page("/create", "Create user page", pageHandler)

		products := users.Group("/products")
		products.AccessGroups("product_manager")
		products.Page("/list", "List products page", pageHandler)

		tests := []struct {
			path           string
			expectedGroups []string
		}{
			{"/users/list", []string{"admin"}},
			{"/users/create", []string{"admin", "customer_support"}},
			{"/users/products/list", []string{"admin", "customer_support", "product_manager"}},
		}

		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				page := findPageByPath(st.pages, tt.path)
				if page == nil {
					t.Fatalf("Page not found: %s", tt.path)
				}

				if len(page.accessGroups) != len(tt.expectedGroups) {
					t.Errorf("Expected %d access groups, got %d for path %s",
						len(tt.expectedGroups), len(page.accessGroups), tt.path)
					t.Errorf("Expected groups: %v, got: %v", tt.expectedGroups, page.accessGroups)
					return
				}

				for _, expectedGroup := range tt.expectedGroups {
					found := false
					for _, actualGroup := range page.accessGroups {
						if actualGroup == expectedGroup {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected group %s not found in %v for path %s",
							expectedGroup, page.accessGroups, tt.path)
					}
				}
			})
		}
	})

	t.Run("Complex group structure", func(t *testing.T) {
		config := &Config{
			APIKey: "test_api_key",
			Host:   "ws://test.trysourcetool.com",
		}
		st := New(config)

		admin := st.Group("/admin")
		admin.AccessGroups("admin")
		admin.Page("/dashboard", "Admin Dashboard", pageHandler)

		settings := admin.Group("/settings")
		settings.AccessGroups("super_admin")
		settings.Page("/system", "System Settings", pageHandler)

		api := st.Group("/api")
		api.AccessGroups("api_user")

		v1 := api.Group("/v1")
		v1.Page("/users", "Users API v1", pageHandler)

		v2 := api.Group("/v2")
		v2.AccessGroups("api_v2")
		v2.Page("/users", "Users API v2", pageHandler)

		tests := []struct {
			path           string
			expectedGroups []string
		}{
			{"/admin/dashboard", []string{"admin"}},
			{"/admin/settings/system", []string{"admin", "super_admin"}},
			{"/api/v1/users", []string{"api_user"}},
			{"/api/v2/users", []string{"api_user", "api_v2"}},
		}

		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				page := findPageByPath(st.pages, tt.path)
				if page == nil {
					t.Fatalf("Page not found: %s", tt.path)
				}

				if len(page.accessGroups) != len(tt.expectedGroups) {
					t.Errorf("Expected %d access groups, got %d for path %s",
						len(tt.expectedGroups), len(page.accessGroups), tt.path)
				}

				for _, expectedGroup := range tt.expectedGroups {
					found := false
					for _, actualGroup := range page.accessGroups {
						if actualGroup == expectedGroup {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected group %s not found in %v for path %s",
							expectedGroup, page.accessGroups, tt.path)
					}
				}
			})
		}
	})

	t.Run("Error handling", func(t *testing.T) {
		config := &Config{
			APIKey: "test_api_key",
			Host:   "ws://test.trysourcetool.com",
		}
		st := New(config)
		errorHandler := func(ui UIBuilder) error {
			return errors.New("test error")
		}
		st.Page("/error", "Error Page", errorHandler)

		page := findPageByPath(st.pages, "/error")
		if page == nil {
			t.Fatal("Page not found")
		}

		if err := page.run(nil); err == nil {
			t.Error("Expected error from handler, got nil")
		}
	})
}

func findPageByPath(pages map[uuid.UUID]*page, path string) *page {
	for _, p := range pages {
		if p.route == path {
			return p
		}
	}
	return nil
}

func TestPageManager(t *testing.T) {
	t.Run("Get existing page", func(t *testing.T) {
		pages := make(map[uuid.UUID]*page)
		pageID := uuid.Must(uuid.NewV4())
		testPage := &page{
			id:   pageID,
			name: "TestPage",
		}
		pages[pageID] = testPage

		pm := newPageManager(pages)

		got := pm.getPage(pageID)
		if got != testPage {
			t.Error("getPage returned wrong page")
		}
	})

	t.Run("Get non-existent page", func(t *testing.T) {
		pages := make(map[uuid.UUID]*page)
		pm := newPageManager(pages)

		nonExistentID := uuid.Must(uuid.NewV4())
		got := pm.getPage(nonExistentID)
		if got != nil {
			t.Error("getPage returned non-nil for non-existent page")
		}
	})
}
