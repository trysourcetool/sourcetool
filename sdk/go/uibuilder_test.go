package sourcetool

import (
	"context"
	"testing"
)

func TestCursor_PathManagement(t *testing.T) {
	c := newCursor()

	// Test initial path
	initialPath := c.getPath()
	if len(initialPath) != 1 || initialPath[0] != 0 {
		t.Errorf("initial path = %v, want [0]", initialPath)
	}

	// Test next()
	c.next()
	nextPath := c.getPath()
	if len(nextPath) != 1 || nextPath[0] != 1 {
		t.Errorf("path after next() = %v, want [1]", nextPath)
	}

	// Add parent path
	c.parentPath = append(c.parentPath, 1)
	parentPath := c.getPath()
	if len(parentPath) != 2 || parentPath[0] != 1 || parentPath[1] != 1 {
		t.Errorf("path with parent = %v, want [1,1]", parentPath)
	}
}

func TestPath_String(t *testing.T) {
	tests := []struct {
		path path
		want string
	}{
		{path{0}, "0"},
		{path{1, 2, 3}, "123"},
		{path{0, 1, 0}, "010"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.path.String(); got != tt.want {
				t.Errorf("path.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUIBuilder_Context(t *testing.T) {
	ctx := context.Background()
	builder := &uiBuilder{
		context: ctx,
	}

	if got := builder.Context(); got != ctx {
		t.Errorf("Context() = %v, want %v", got, ctx)
	}
}
