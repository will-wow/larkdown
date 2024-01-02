package larkdown

import (
	"errors"

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
	opts ...FindOption,
) (out T, err error) {
	config := newFindConfig(opts...)

	found, err := query.QueryOne(doc, source, matcher)
	if err != nil {
		// Return nil for a no match error if allowed.
		var queryErr *query.QueryError
		if config.AllowNoMatch && errors.As(err, &queryErr) {
			var empty T
			return empty, nil
		}

		return out, err
	}

	return fn(found, source)
}

// FindConfig configures the Find function.
type FindConfig struct {
	// AllowNoMatch allows Find to return nil when no match is found. By default it will return a query.QueryError.
	AllowNoMatch bool
}

// FindOption describes a functional option for Find.
type FindOption func(*FindConfig)

// newFindConfig returns a new FindConfig with default values.
func newFindConfig(opts ...FindOption) *FindConfig {
	config := &FindConfig{
		AllowNoMatch: false,
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

// FindAllowNoMatch allows Find to return nil when no match is found. By default it will return a query.QueryError.
func FindAllowNoMatch() FindOption {
	return func(config *FindConfig) {
		config.AllowNoMatch = true
	}
}

// Use a matcher to find a all nodes, then decode its contents and return structured data.
func FindAll[T any](
	doc ast.Node,
	source []byte,
	matcher []match.Node,
	extractor match.Node,
	fn func(node ast.Node, source []byte) (T, error),
	opts ...FindAllOption,
) (out []T, err error) {
	config := newFindAllConfig(opts...)

	found, err := query.QueryAll(doc, source, matcher, extractor)
	if err != nil {
		// Return nil for a no match error if allowed.
		var queryErr *query.QueryError
		if config.AllowNoMatch && errors.As(err, &queryErr) {
			return nil, nil
		}

		return nil, err
	}

	out = make([]T, len(found))

	for i, node := range found {
		decoded, err := fn(node, source)
		if err != nil {
			return out, err
		}

		out[i] = decoded
	}

	return out, nil
}

// FindAllConfig configures the FindAll function.
type FindAllConfig struct {
	// AllowNoMatch allows FindAll to return nil when no match is found. By default it will return a query.QueryError.
	AllowNoMatch bool
}

// FindAllOption describes a functional option for FindAll.
type FindAllOption func(*FindAllConfig)

// newFindAllConfig returns a new FindAllConfig with default values.
func newFindAllConfig(opts ...FindAllOption) *FindAllConfig {
	config := &FindAllConfig{
		AllowNoMatch: false,
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

// AllowNoMatch allows FindAll to return nil when no match is found. By default it will return a query.QueryError.
func FindAllAllowNoMatch() FindAllOption {
	return func(config *FindAllConfig) {
		config.AllowNoMatch = true
	}
}

// NewNodeRenderer returns a new goldmark NodeRenderer with default config that renders nodes as Markdown.
func NewNodeRenderer(opts ...mdrender.Option) renderer.Renderer {
	return renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(mdrender.NewRenderer(opts...), 998)))
}
