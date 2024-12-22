// Package iterutil provides iterator utility functions.
package iterutil

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

// Ordered returns an iterator over key-value pairs from m.
// The iteration order is deterministic, arranging keys in ascending order.
func Ordered[Map ~map[K]V, K cmp.Ordered, V any](m Map) iter.Seq2[K, V] {
	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)
	return func(yield func(K, V) bool) {
		for _, k := range keys {
			v := m[k]
			if !yield(k, v) {
				return
			}
		}
	}
}
