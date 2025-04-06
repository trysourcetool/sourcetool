package sourcetool

import (
	"errors"
	"testing"
)

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		path     string
		want     string
	}{
		{
			name:     "Empty base path",
			basePath: "",
			path:     "/users",
			want:     "/users",
		},
		{
			name:     "Base path with trailing slash",
			basePath: "/admin/",
			path:     "users",
			want:     "/admin/users",
		},
		{
			name:     "Path without leading slash",
			basePath: "/admin",
			path:     "users",
			want:     "/admin/users",
		},
		{
			name:     "Both with slashes",
			basePath: "/admin/",
			path:     "/users/",
			want:     "/admin/users",
		},
		{
			name:     "Nested paths",
			basePath: "/api/v1",
			path:     "users/list",
			want:     "/api/v1/users/list",
		},
		{
			name:     "Root path",
			basePath: "",
			path:     "/",
			want:     "/",
		},
		{
			name:     "Root path with base path",
			basePath: "/admin",
			path:     "/",
			want:     "/admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router{basePath: tt.basePath}
			got := r.joinPath(tt.path)
			if got != tt.want {
				t.Errorf("joinPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name   string
		groups []string
		want   []string
	}{
		{
			name:   "No duplicates",
			groups: []string{"admin", "user", "guest"},
			want:   []string{"admin", "user", "guest"},
		},
		{
			name:   "With duplicates",
			groups: []string{"admin", "user", "admin", "guest", "user"},
			want:   []string{"admin", "user", "guest"},
		},
		{
			name:   "Empty slice",
			groups: []string{},
			want:   []string{},
		},
		{
			name:   "Single element",
			groups: []string{"admin"},
			want:   []string{"admin"},
		},
		{
			name:   "All duplicates",
			groups: []string{"admin", "admin", "admin"},
			want:   []string{"admin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeDuplicates(tt.groups)
			if len(got) != len(tt.want) {
				t.Errorf("removeDuplicates() length = %v, want %v", len(got), len(tt.want))
			}
			for _, w := range tt.want {
				found := false
				for _, g := range got {
					if g == w {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("removeDuplicates() missing element %v", w)
				}
			}
		})
	}
}

func TestGeneratePageID(t *testing.T) {
	r := &router{
		namespaceDNS: "test.trysourcetool.com",
	}

	tests := []struct {
		name     string
		path     string
		wantSame bool
	}{
		{
			name:     "Simple path",
			path:     "/users",
			wantSame: true,
		},
		{
			name:     "Nested path",
			path:     "/admin/users/list",
			wantSame: true,
		},
		{
			name:     "Root path",
			path:     "/",
			wantSame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id1 := r.generatePageID(tt.path)
			id2 := r.generatePageID(tt.path)

			if tt.wantSame && id1 != id2 {
				t.Errorf("generatePageID() generated different IDs for same path")
			}

			differentPath := tt.path + "/different"
			id3 := r.generatePageID(differentPath)
			if id1 == id3 {
				t.Errorf("generatePageID() generated same ID for different paths")
			}
		})
	}
}

func TestRouterAccessGroups(t *testing.T) {
	pageHandler := func(ui UIBuilder) error { return nil }

	t.Run("Group creation before and after AccessGroups", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)

		st.AccessGroups("global")
		admin := st.Group("/admin")
		admin.Page("/dashboard", "Dashboard", pageHandler)

		users := st.Group("/users")
		users.AccessGroups("user_manager")
		users.Page("/list", "User List", pageHandler)

		tests := []struct {
			path           string
			expectedGroups []string
		}{
			{"/admin/dashboard", []string{"global"}},
			{"/users/list", []string{"global", "user_manager"}},
		}

		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				page := findPageByPath(st.pages, tt.path)
				assertPageGroups(t, page, tt.expectedGroups)
			})
		}
	})

	t.Run("Multiple AccessGroups calls", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		admin := st.Group("/admin")

		admin.AccessGroups("admin")
		admin.AccessGroups("super_admin")
		admin.Page("/settings", "Settings", pageHandler)

		userAdminGroup := admin.Group("/")
		userAdminGroup.AccessGroups("user_admin").Page("/users", "Users", pageHandler)
		systemAdminGroup := admin.Group("/")
		systemAdminGroup.AccessGroups("system_admin").Page("/system", "System", pageHandler)

		tests := []struct {
			path           string
			expectedGroups []string
		}{
			{"/admin/settings", []string{"admin", "super_admin"}},
			{"/admin/users", []string{"admin", "super_admin", "user_admin"}},
			{"/admin/system", []string{"admin", "super_admin", "system_admin"}},
		}

		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				page := findPageByPath(st.pages, tt.path)
				assertPageGroups(t, page, tt.expectedGroups)
			})
		}
	})

	t.Run("Sibling groups inheritance", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		st.AccessGroups("global")

		users := st.Group("/users")
		users.AccessGroups("user_admin")
		users.Page("/list", "Users", pageHandler)

		products := st.Group("/products")
		products.AccessGroups("product_admin")
		products.Page("/list", "Products", pageHandler)

		tests := []struct {
			path           string
			expectedGroups []string
		}{
			{"/users/list", []string{"global", "user_admin"}},
			{"/products/list", []string{"global", "product_admin"}},
		}

		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				page := findPageByPath(st.pages, tt.path)
				assertPageGroups(t, page, tt.expectedGroups)
			})
		}
	})

	t.Run("Deep nested groups inheritance", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		st.AccessGroups("global")

		api := st.Group("/api")
		api.AccessGroups("api_user")

		v1 := api.Group("/v1")
		v1.AccessGroups("v1_user")

		users := v1.Group("/users")
		users.AccessGroups("user_admin")

		settings := users.Group("/settings")
		settings.AccessGroups("settings_admin")
		settings.Page("/profile", "Profile Settings", pageHandler)

		page := findPageByPath(st.pages, "/api/v1/users/settings/profile")
		expectedGroups := []string{"global", "api_user", "v1_user", "user_admin", "settings_admin"}
		assertPageGroups(t, page, expectedGroups)
	})

	t.Run("Mixed group and page specific access groups", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)

		admin := st.Group("/admin")
		admin.AccessGroups("admin")
		admin.Page("/dashboard", "Dashboard", pageHandler)

		settings := admin.Group("/settings")
		settings.AccessGroups("settings_admin")
		settings.AccessGroups("system_admin")
		settings.Page("/general", "General Settings", pageHandler)

		users := settings.Group("/users")
		users.AccessGroups("user_manager")
		users.AccessGroups("profile_admin").Page("/profiles", "User Profiles", pageHandler)

		tests := []struct {
			path           string
			expectedGroups []string
		}{
			{"/admin/dashboard", []string{"admin"}},
			{"/admin/settings/general", []string{"admin", "settings_admin", "system_admin"}},
			{"/admin/settings/users/profiles", []string{"admin", "settings_admin", "system_admin", "user_manager", "profile_admin"}},
		}

		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				page := findPageByPath(st.pages, tt.path)
				assertPageGroups(t, page, tt.expectedGroups)
			})
		}
	})
}

func assertPageGroups(t *testing.T, page *page, expectedGroups []string) {
	t.Helper()

	if page == nil {
		t.Fatal("Page not found")
	}

	if len(page.accessGroups) != len(expectedGroups) {
		t.Errorf("Expected %d access groups, got %d", len(expectedGroups), len(page.accessGroups))
		t.Errorf("Expected groups: %v, got: %v", expectedGroups, page.accessGroups)
		return
	}

	for _, expectedGroup := range expectedGroups {
		found := false
		for _, actualGroup := range page.accessGroups {
			if actualGroup == expectedGroup {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected group %s not found in %v", expectedGroup, page.accessGroups)
		}
	}
}

func TestRouterGroup(t *testing.T) {
	pageHandler := func(ui UIBuilder) error { return nil }

	t.Run("Base path construction", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		admin := st.Group("/admin")
		settings := admin.Group("/settings")
		settings.Page("/users", "User Settings", pageHandler)

		page := findPageByPath(st.pages, "/admin/settings/users")
		if page == nil {
			t.Fatal("Page not found")
		}

		if page.route != "/admin/settings/users" {
			t.Errorf("Expected route /admin/settings/users, got %s", page.route)
		}
	})

	t.Run("Multiple nested groups", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		api := st.Group("/api")
		v1 := api.Group("/v1")
		users := v1.Group("/users")
		users.Page("/list", "List Users", pageHandler)

		page := findPageByPath(st.pages, "/api/v1/users/list")
		if page == nil {
			t.Fatal("Page not found")
		}

		if page.route != "/api/v1/users/list" {
			t.Errorf("Expected route /api/v1/users/list, got %s", page.route)
		}
	})
}

func TestRouter_Page(t *testing.T) {
	t.Run("Basic page", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		handler := func(ui UIBuilder) error { return nil }
		st.Page("/test", "Test Page", handler)

		page := findPageByPath(st.pages, "/test")
		if page == nil {
			t.Fatal("Page not found")
		}

		if page.name != "Test Page" {
			t.Errorf("Page name = %v, want %v", page.name, "Test Page")
		}
	})

	t.Run("Skip top-level root path", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		handler := func(ui UIBuilder) error { return nil }
		st.Page("/", "Root Page", handler)

		page := findPageByPath(st.pages, "/")
		if page != nil {
			t.Error("Top-level root path page should not be created")
		}

		// Verify that other pages can still be created
		st.Page("/other", "Other Page", handler)
		otherPage := findPageByPath(st.pages, "/other")
		if otherPage == nil {
			t.Error("Other page should be created")
		}
	})

	t.Run("Allow nested root path", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		handler := func(ui UIBuilder) error { return nil }

		users := st.Group("/users")
		users.Page("/", "Users List", handler)

		page := findPageByPath(st.pages, "/users")
		if page == nil {
			t.Error("Nested root path page should be created")
		}
		if page != nil && page.name != "Users List" {
			t.Errorf("Page name = %v, want %v", page.name, "Users List")
		}

		// Test deeply nested root path
		admin := st.Group("/admin")
		settings := admin.Group("/settings")
		settings.Page("/", "Settings Home", handler)

		settingsPage := findPageByPath(st.pages, "/admin/settings")
		if settingsPage == nil {
			t.Error("Deeply nested root path page should be created")
		}
		if settingsPage != nil && settingsPage.name != "Settings Home" {
			t.Errorf("Page name = %v, want %v", settingsPage.name, "Settings Home")
		}
	})

	t.Run("Page with access groups", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		st.AccessGroups("admin")
		handler := func(ui UIBuilder) error { return nil }
		st.Page("/admin", "Admin Page", handler)

		page := findPageByPath(st.pages, "/admin")
		if page == nil {
			t.Fatal("Page not found")
		}

		if len(page.accessGroups) != 1 || page.accessGroups[0] != "admin" {
			t.Errorf("Expected [admin] access group, got %v", page.accessGroups)
		}
	})

	t.Run("Page with error handler", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		handler := func(ui UIBuilder) error {
			return errors.New("test error")
		}
		st.Page("/error", "Error Page", handler)

		page := findPageByPath(st.pages, "/error")
		if page == nil {
			t.Fatal("Page not found")
		}

		if err := page.run(nil); err == nil {
			t.Error("Expected error from handler, got nil")
		}
	})

	t.Run("Page with empty route", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		handler := func(ui UIBuilder) error { return nil }
		st.Page("", "Root Page", handler)

		page := findPageByPath(st.pages, "/")
		if page == nil {
			t.Fatal("Page not found")
		}

		if page.route != "/" {
			t.Errorf("Expected route '/', got %q", page.route)
		}

		if page.name != "Root Page" {
			t.Errorf("Expected page name 'Root Page', got %q", page.name)
		}
	})

	t.Run("Page with duplicate route", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		handler := func(ui UIBuilder) error { return nil }
		st.Page("/duplicate", "First Page", handler)
		st.Page("/duplicate", "Second Page", handler)

		page := findPageByPath(st.pages, "/duplicate")
		if page == nil {
			t.Fatal("Page not found")
		}

		if page.name != "Second Page" {
			t.Errorf("Expected page name to be 'Second Page', got %q", page.name)
		}
	})
}

func TestRouter_Group(t *testing.T) {
	t.Run("Basic group", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		group := st.Group("/test")

		if group == nil {
			t.Fatal("Group returned nil")
		}

		handler := func(ui UIBuilder) error { return nil }
		group.Page("/page", "Test Page", handler)

		page := findPageByPath(st.pages, "/test/page")
		if page == nil {
			t.Fatal("Page not found")
		}
	})

	t.Run("Group with access groups", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		group := st.Group("/admin")
		group.AccessGroups("admin")

		handler := func(ui UIBuilder) error { return nil }
		group.Page("/dashboard", "Admin Dashboard", handler)

		page := findPageByPath(st.pages, "/admin/dashboard")
		if page == nil {
			t.Fatal("Page not found")
		}

		if len(page.accessGroups) != 1 || page.accessGroups[0] != "admin" {
			t.Errorf("Expected [admin] access group, got %v", page.accessGroups)
		}
	})

	t.Run("Nested groups", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		parent := st.Group("/parent")
		child := parent.Group("/child")

		if child == nil {
			t.Fatal("Child group returned nil")
		}

		handler := func(ui UIBuilder) error { return nil }
		child.Page("/page", "Test Page", handler)

		page := findPageByPath(st.pages, "/parent/child/page")
		if page == nil {
			t.Fatal("Page not found")
		}
	})

	t.Run("Group with empty path", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		group := st.Group("")

		if group == nil {
			t.Fatal("Group returned nil")
		}

		handler := func(ui UIBuilder) error { return nil }
		group.Page("/page", "Test Page", handler)

		page := findPageByPath(st.pages, "/page")
		if page == nil {
			t.Fatal("Page not found")
		}
	})
}

func TestRouter_AccessGroups(t *testing.T) {
	t.Run("Set access groups", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		st.AccessGroups("admin", "user")

		handler := func(ui UIBuilder) error { return nil }
		st.Page("/test", "Test Page", handler)

		page := findPageByPath(st.pages, "/test")
		if page == nil {
			t.Fatal("Page not found")
		}

		if len(page.accessGroups) != 2 {
			t.Errorf("Expected 2 access groups, got %d", len(page.accessGroups))
		}

		expectedGroups := []string{"admin", "user"}
		for _, group := range expectedGroups {
			found := false
			for _, actualGroup := range page.accessGroups {
				if actualGroup == group {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected group %s not found in %v", group, page.accessGroups)
			}
		}
	})

	t.Run("Clear access groups", func(t *testing.T) {
		config := &Config{
			APIKey:   "test_api_key",
			Endpoint: "ws://test.trysourcetool.com",
		}
		st := New(config)
		st.AccessGroups("admin")
		st.AccessGroups()

		handler := func(ui UIBuilder) error { return nil }
		st.Page("/test", "Test Page", handler)

		page := findPageByPath(st.pages, "/test")
		if page == nil {
			t.Fatal("Page not found")
		}

		if len(page.accessGroups) != 1 || page.accessGroups[0] != "admin" {
			t.Errorf("Expected [admin] access group, got %v", page.accessGroups)
		}
	})
}
