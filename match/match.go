// match provides a query language for matching nodes in a larkdown.Tree
package match

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown/gmast"
)

// Interface for a node matcher.
type Node interface {
	// Returns true if the node matches the query.
	Match(node ast.Node, index int, source []byte) (ok bool)
	// Returns true if the node is the end of a flat branch.
	// This is used by headings to note when another heading ends the branch.
	EndMatch(node ast.Node) bool
	// Show the matcher as a string in error messages.
	String() string
	// Returns the next node to match, which usually a child, or for headings a sibling.
	NextNode(self ast.Node) ast.Node
	// True for nodes whose children are actually siblings.
	IsFlatBranch() bool
}

// Partially implements Node.
type BaseNode struct{}

func (m BaseNode) EndMatch(node ast.Node) bool {
	return false
}
func (m BaseNode) NextNode(self ast.Node) ast.Node {
	return self.FirstChild()
}
func (m BaseNode) IsFlatBranch() bool {
	return false
}

// Matches a heading by level and name.
type Branch struct {
	// The heading level to match, or 0 to match any level.
	Level int
	// The heading name to match, or empty to match any name.
	Name []byte
	// If true, the name is matched case-insensitively.
	CaseInsensitive bool
}

var _ Node = Branch{}

// Match a heading by level and name.
func (m Branch) Match(node ast.Node, index int, source []byte) bool {
	heading, ok := node.(*ast.Heading)
	if !ok {
		return false
	}
	if m.Level != 0 && heading.Level != m.Level {
		return false
	}

	// If the name is empty, we're matching any heading of the given level.
	if len(m.Name) == 0 {
		return true
	}

	if m.CaseInsensitive {
		return bytes.EqualFold(node.FirstChild().Text(source), m.Name)
	}
	return bytes.Equal(node.FirstChild().Text(source), m.Name)
}

// A heading branch ends when the next heading is of the same or higher level.
func (m Branch) EndMatch(node ast.Node) bool {
	heading, ok := node.(*ast.Heading)
	if !ok {
		return false
	}

	return heading.Level <= m.Level
}

// For headings, the next node is the next sibling.
func (m Branch) NextNode(self ast.Node) ast.Node {
	return gmast.GetNextSibling(self)
}

// Headings are branches whose children are siblings until the next heading.
func (m Branch) IsFlatBranch() bool {
	return true
}

func (m Branch) String() string {
	var level string
	if m.Level == 0 {
		// If the level is unspecified, note that.
		level = "#?"
	} else {
		// Otherwise indicate the level with hashes.
		level = strings.Repeat("#", m.Level)
	}

	// If the name is empty, we're matching any heading of the given level.
	if len(m.Name) == 0 {
		return fmt.Sprintf("[%s]", level)
	}

	return fmt.Sprintf("[%s %s]", level, string(m.Name))
}

// Matches an ordered or unordered list.
type List struct {
	BaseNode
}

var _ Node = List{}

// Matches List nodes.
func (m List) Match(node ast.Node, index int, source []byte) bool {
	_, ok := node.(*ast.List)
	return ok
}

func (m List) String() string {
	return ".list"
}

// Matches go.abhg.dev/goldmark/hashtag #tag nodes.
type Tag struct {
	BaseNode
}

var _ Node = Tag{}

func (m Tag) Match(node ast.Node, index int, source []byte) bool {
	_, ok := node.(*hashtag.Node)
	return ok
}

func (m Tag) String() string {
	return "[#tag]"
}

// Wraps another query, only when it's the nth child of the parent.
type Index struct {
	Index int
	Node  Node
}

var _ Node = Index{}

func (m Index) Match(node ast.Node, index int, source []byte) bool {
	if m.Index != index {
		return false
	}

	return m.Node.Match(node, index, source)
}

func (m Index) EndMatch(node ast.Node) bool {
	return m.Node.EndMatch(node)
}

func (m Index) NextNode(self ast.Node) ast.Node {
	return m.Node.NextNode(self)
}

func (m Index) IsFlatBranch() bool {
	return m.Node.IsFlatBranch()
}

func (m Index) String() string {
	return fmt.Sprintf("[%d]%s", m.Index, m.Node.String())
}

// Matches against a specific node kind.
// Used for matching arbitrary nodes, including custom ones.
type NodeOfKind struct {
	BaseNode
	Kind ast.NodeKind
}

var _ Node = NodeOfKind{}

func (m NodeOfKind) Match(node ast.Node, index int, source []byte) bool {
	return node.Kind() == m.Kind
}

func (m NodeOfKind) String() string {
	return fmt.Sprintf("[kind:%s]", m.Kind.String())
}

// Matches any node. Useful as a fallback for index matches.
type AnyNode struct {
	BaseNode
}

var _ Node = AnyNode{}

func (m AnyNode) Match(node ast.Node, index int, source []byte) bool {
	return true
}

func (m AnyNode) String() string {
	return ".any"
}
