package stream

import "github.com/mt-inside/dagger/pkg/value"

type updateTarget interface {
	Update(value.Value)
}
