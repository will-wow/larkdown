package gmast

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"
)

// AppendChild is a chainable version of ast.Node.AppendChild.
// Builds a new AST in a single call.
func AppendChild[T ast.Node](parent T, children ...ast.Node) (theParent T) {
	for _, child := range children {
		parent.AppendChild(parent, child)
	}
	return parent
}

// AppendHeadingChild appends a child to the last sibling of the heading
// before a new higher-level heading. Conceptually, this is the same as AppendChild.
func AppendHeadingChild[T ast.Node](heading T, children ...ast.Node) (theParent T) {
	lastChild, err := LastChildOfHeading(heading)
	if err != nil {
		return heading
	}

	parent := heading.Parent()
	if parent == nil {
		return heading
	}

	for _, child := range children {
		parent.InsertAfter(parent, lastChild, child)
	}

	return heading
}

// NewTextSegment builds a new ast.TextSegment and appends its content to the source in a single call.
func NewTextSegment(newText string, source []byte) (textSegment *ast.Text, newSource []byte) {
	newSegment, newSource := NewSegment(newText, source)

	textSegment = ast.NewTextSegment(newSegment)

	return textSegment, newSource
}

// NewSegment builds a new text.Segment and appends its content to the source in a single call.
func NewSegment(newText string, source []byte) (text.Segment, []byte) {
	newSegment := text.NewSegment(len(source), len(source)+len(newText))

	source = append(source, newText...)

	return newSegment, source
}

// NewParagraph is a helper to build a new paragraph with a single text segment.
func NewParagraph(newText string, source []byte) (node ast.Node, newSource []byte) {
	newSegment, newSource := NewSegment(newText, source)

	p := AppendChild(
		ast.NewParagraph(),
		AppendChild(
			ast.NewTextBlock(),
			ast.NewTextSegment(newSegment),
		))

	return p, newSource
}

// NewSpace returns a new space text segment for adding space between inline nodes.
func NewSpace(source []byte) (node ast.Node, newSource []byte) {
	return NewTextSegment(" ", source)
}

// NewHashtag returns a new hashtag.Node with the text set up correctly.
func NewHashtag(tag string, source []byte) (node ast.Node, newSource []byte) {
	taggedTag := fmt.Sprintf("#%s", tag)
	newText, newSource := NewTextSegment(taggedTag, source)

	return AppendChild(
			// A tag with the tag
			&hashtag.Node{Tag: []byte(tag)},
			// The child of the hashtag must be #tag to render correctly
			newText,
		),
		newSource
}
