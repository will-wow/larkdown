// Package mdfront adds support for rendering frontmatter to markdown for goldmark.
package mdfront

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Kind is the kind of hashtag AST nodes.
var Kind = ast.NewNodeKind("Frontmatter")

// Node is a hashtag node in a Goldmark Markdown document.
type Node struct {
	ast.BaseInline

	// Frontmatter is a struct with yaml tags for rendering.
	Frontmatter any
}

var _ ast.Node = &Node{}

// Kind reports the kind of hashtag nodes.
func (*Node) Kind() ast.NodeKind { return Kind }

// Kind reports the kind of hashtag nodes.
func (*Node) IsRaw() bool { return true }

// Dump dumps the contents of Node to stdout for debugging.
func (n *Node) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"Frontmatter": fmt.Sprintf("%+v", n.Frontmatter),
	}, nil)
}

// astTransformer adds a blank frontmatter node for filling with data during markdown rendering.
type astTransformer struct {
}

// Transform adds a new frontmatter node to the AST.
func (a *astTransformer) Transform(doc *gast.Document, reader text.Reader, pc parser.Context) {
	fontMatterNode := &Node{}

	doc.InsertBefore(doc, doc.FirstChild(), fontMatterNode)
}

// NewTransformer adds a blank frontmatter node for filling with data during markdown rendering.
func NewTransformer() parser.ASTTransformer {
	return &astTransformer{}
}

// NoopRenderer renders nothing for a frontmatter block,
// for if the AST is rendered without the mdrender.Renderer
type NoopRenderer struct{}

var _ renderer.NodeRenderer = &NoopRenderer{}

func (r *NoopRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, renderNothing)
}

// renderNothing renders nothing. Instead mdrender.Renderer.renderNothing
// will render to markdown because it is higher priority.
func renderNothing(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

// Extender extends a goldmark Markdown object with support for setting
// up a frontmatter placeholder for later rendering.
//
// Install it on your Markdown object upon creation.
//
//	goldmark.New(
//	  goldmark.WithExtensions(
//	    // ...
//	    &mdfront.Extender{},
//	  ),
//	  // ...
//	)
type Extender struct{}

var _ goldmark.Extender = (*Extender)(nil)

// Extend extends the provided goldmark Markdown object with support for
// a placeholder for rendering frontmatter.
func (e *Extender) Extend(m goldmark.Markdown) {
	// Adds a frontmatter node to front the AST, to be rendered.
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(NewTransformer(), 0),
		),
	)
	// Add a catch-all to not render the frontmatter node
	// with the default renderer.
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&NoopRenderer{}, 10_000),
		),
	)
}
