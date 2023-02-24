// Package pointer provides useful functions to convert between pointers and values.
package pointer

// New returns pointer of given variable.
func New[T any](v T) *T {
	return &v
}

// From returns value of given pointer. Returns default value if pointer is nil.
func From[T any](p *T) T {
	if p == nil {
		var v T
		return v
	}

	return *p
}

// NewSlice returns a slice of pointers given a slice of values.
func NewSlice[T any](in []T) []*T {
	out := make([]*T, len(in))
	for i := range in {
		out[i] = New(in[i])
	}
	return out
}

// FromSlice returns a slice of values given a slice of pointers. Ignores nil elements.
func FromSlice[T any](in []*T) []T {
	out := make([]T, 0, len(in))
	for i := range in {
		if in[i] == nil {
			continue
		}

		out = append(out, *in[i])
	}
	return out
}
