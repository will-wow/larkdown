package larkdown

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"

	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/mdrender"
	"github.com/will-wow/larkdown/query"
)

// Use a matcher to find a node, and then decode its contents and return structured data.
func Find[T any](
	doc ast.Node,
	source []byte,
	matcher []match.Node,
	fn func(node ast.Node, source []byte) (T, error),
) (out T, err error) {
	found, err := query.QueryOne(doc, source, matcher)
	if err != nil {
		return out, err
	}

	return fn(found, source)
}

// Use a matcher to find a all nodes, then decode its contents and return structured data.
func FindAll[T any](
	doc ast.Node,
	source []byte,
	matcher []match.Node,
	extractor match.Node,
	fn func(node ast.Node, source []byte) (T, error),
) (out []T, err error) {
	found, err := query.QueryAll(doc, source, matcher, extractor)
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

// NewNodeRenderer returns a new goldmark NodeRenderer with default config that renders nodes as Markdown.
func NewNodeRenderer(opts ...mdrender.Option) renderer.Renderer {
	return renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(mdrender.NewRenderer(opts...), 998)))
}
