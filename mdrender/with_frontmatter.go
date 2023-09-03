package mdrender

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Kind is the kind of hashtag AST nodes.
var FrontmatterKind = ast.NewNodeKind("Frontmatter")

// Node is a hashtag node in a Goldmark Markdown document.
type FrontmatterNode struct {
	ast.BaseInline

	// Frontmatter is a struct with yaml tags for rendering.
	Frontmatter any
}

var _ ast.Node = &FrontmatterNode{}

// Kind reports the kind of hashtag nodes.
func (*FrontmatterNode) Kind() ast.NodeKind { return FrontmatterKind }

// Kind reports the kind of hashtag nodes.
func (*FrontmatterNode) IsRaw() bool { return true }

// Dump dumps the contents of Node to stdout for debugging.
func (n *FrontmatterNode) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"Frontmatter": fmt.Sprintf("%+v", n.Frontmatter),
	}, nil)
}

// func NewFrontmatterNode() *FrontmatterNode {
// 	newSegment, newSource := gmast.NewSegment("", source)
// 	return gmast.NewParagraph() & FrontmatterNode{}
// }

// astTransformer adds a blank frontmatter node for filling with data during markdown rendering.
type astTransformer struct {
}

func (a *astTransformer) Transform(doc *gast.Document, reader text.Reader, pc parser.Context) {
	fontMatterNode := &FrontmatterNode{}

	doc.InsertBefore(doc, doc.FirstChild(), fontMatterNode)
}

// NewTransformer adds a blank frontmatter node for filling with data during markdown rendering.
func NewTransformer() parser.ASTTransformer {
	return &astTransformer{}
}
