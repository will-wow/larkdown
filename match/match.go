// match provides a query language for matching nodes in a larkdown.Tree
package match

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
)

// Interface for a node matcher.
type Node interface {
	Match(node ast.Node, index int, source []byte) (ok bool)
	EndMatch(node ast.Node, index int, source []byte) bool
	String() string
	ShouldDrill() bool
	IsFlatBranch() bool
}

// Matches a heading by level and name.
type Branch struct {
	Level           int
	Name            []byte
	CaseInsensitive bool
}

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

func (m Branch) EndMatch(node ast.Node, index int, source []byte) bool {
	heading, ok := node.(*ast.Heading)
	if !ok {
		return false
	}

	return heading.Level <= m.Level
}

func (m Branch) ShouldDrill() bool {
	return false
}

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
type List struct{}

func (m List) Match(node ast.Node, index int, source []byte) bool {
	_, ok := node.(*ast.List)
	return ok
}

func (m List) EndMatch(node ast.Node, index int, source []byte) bool {
	return true
}

func (m List) ShouldDrill() bool {
	return true
}

func (m List) IsFlatBranch() bool {
	return false
}

func (m List) String() string {
	return ".list"
}

// Wraps another query, only when it's the nth child of the parent.
type Index struct {
	Index int
	Node  Node
}

func NewIndex(index int, node Node) *Index {
	return &Index{
		Index: index,
		Node:  node,
	}
}

func (m Index) Match(node ast.Node, index int, source []byte) bool {
	if m.Index != index {
		return false
	}

	return m.Node.Match(node, index, source)
}

func (m Index) EndMatch(node ast.Node, index int, source []byte) bool {
	return m.Node.EndMatch(node, index, source)
}

func (m Index) ShouldDrill() bool {
	return m.Node.ShouldDrill()
}

func (m Index) IsFlatBranch() bool {
	return m.Node.IsFlatBranch()
}

func (m Index) String() string {
	return fmt.Sprintf("[%d]%s", m.Index, m.Node.String())
}

type SearchFor struct {
	Node Node
}

func NewSearchFor() *SearchFor {
	return &SearchFor{}
}

func (m SearchFor) Match(node ast.Node, index int, source []byte) bool {
	return true
}

func (m SearchFor) String() string {
	return ".searchFor"
}

func (m SearchFor) EndMatch(node ast.Node, index int, source []byte) bool {
	return false
}

func (m SearchFor) ShouldDrill() bool {
	return true
}

func (m SearchFor) IsFlatBranch() bool {
	return false
}

// Matches any node. Useful as a fallback for index matches.
type AnyNode struct{}

func NewAnyNode() *AnyNode {
	return &AnyNode{}
}

func (m AnyNode) Match(node ast.Node, index int, source []byte) bool {
	return true
}

func (m AnyNode) String() string {
	return ".any"
}

func (m AnyNode) EndMatch(node ast.Node, index int, source []byte) bool {
	return false
}

func (m AnyNode) ShouldDrill() bool {
	return true
}

func (m AnyNode) IsFlatBranch() bool {
	return false
}
