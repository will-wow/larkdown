// Package match provides a query language for matching nodes in a larkdown.Tree
package match

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/preprocess"
)

// Interface for a node matcher.
type Node interface {
	Match(node ast.Node, index int, source []byte) (ok bool)
	String() string
}

// Matches a heading by level and name.
type Branch struct {
	Level           int
	Name            []byte
	CaseInsensitive bool
}

func (m Branch) Match(node ast.Node, index int, source []byte) bool {
	heading, ok := node.(*preprocess.TreeBranch)
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

func (m List) String() string {
	return ".list"
}

// Wraps another query, only when it's the nth child of the parent.
// Note that for branches, the heading itself is the 0th child.
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

func (m Index) String() string {
	return fmt.Sprintf("[%d]%s", m.Index, m.Node.String())
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
