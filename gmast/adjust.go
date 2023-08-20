package gmast

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// AppendChild is a chainable version of ast.Node.AppendChild.
// Builds a new AST in a single call.
func AppendChild[T ast.Node](parent T, child ast.Node) (theParent T) {
	parent.AppendChild(parent, child)
	return parent
}

// NewSegment builds a new text.Segment and appends its content to the source in a single call.
func NewSegment(newText string, source []byte) (text.Segment, []byte) {
	newSegment := text.NewSegment(len(source), len(source)+len(newText))

	source = append(source, newText...)

	return newSegment, source
}

func Test() {
}
