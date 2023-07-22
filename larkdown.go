package larkdown

import (
	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
)

// Use a matcher to find a node, and then unmarshal its contents into structured data.
func Find[T any](doc ast.Node, source []byte, matcher []match.Node, fn func(node ast.Node, source []byte) (T, error)) (out T, err error) {
	found, err := query.QueryTree(doc, source, matcher)
	if err != nil {
		return out, err
	}

	return fn(found, source)
}

func FindAll[T any](doc ast.Node, source []byte, matcher []match.Node, fn func(node ast.Node, source []byte) (T, error)) (out []T, err error) {
	found, err := query.QueryAll(doc, source, matcher)
	if err != nil {
		return
	}

	out = make([]T, len(found))

	for i, node := range found {
		decoded, err := fn(node, source)
		if err != nil {
			return out, err
		}

		out[i] = decoded
	}

	return
}
