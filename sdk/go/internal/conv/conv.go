package conv

import (
	"github.com/samber/lo"
)

func SafeValue[T comparable](in *T) T {
	if in == nil {
		return lo.Empty[T]()
	}
	return *in
}

func NilValue[T comparable](in T) *T {
	return &in
}
