package mdrender

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// noopRenderer renders nothing for a frontmatter block,
// for if the AST is rendered without the mdrender.Renderer
type noopRenderer struct{}

var _ renderer.NodeRenderer = &Renderer{}

func (r *noopRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(FrontmatterKind, r.renderFrontmatter)
}

// renderFrontmatter renders nothing by default. mdrender.Renderer.renderFrontmatter
// will render to markdown because it is higher priority.
func (r *noopRenderer) renderFrontmatter(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

// Extender extends a goldmark Markdown object with support for parsing and
// rendering frontmatter to markdown.
//
// Install it on your Markdown object upon creation.
//
//	goldmark.New(
//	  goldmark.WithExtensions(
//	    // ...
//	    &mdrender.Extender{},
//	  ),
//	  // ...
//	)
type Extender struct{}

var _ goldmark.Extender = (*Extender)(nil)

// Extend extends the provided goldmark Markdown object with support for
// a placeholder for rendering frontmatter.
func (e *Extender) Extend(m goldmark.Markdown) {
	//
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(NewTransformer(), 0),
		),
	)
	// Add a catch-all to not render the frontmatter node
	// with the default renderer.
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&noopRenderer{}, 10_000),
		),
	)
}
