// Package match provides a query language for matching nodes in a larkdown.Tree
package match

import (
	"bytes"

	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/preprocess"
)

// TODO: queries:
// Code block with language
// Nth Instance Query that retunrns a "don't pop this query", and keeps state of how many times it's called?

// Interface for a node matcher.
type Node interface {
	Match(node ast.Node, index int, source []byte) (ok bool)
}

// Matches a heading by level and name.
type Branch struct {
	Level           int
	Name            []byte
	CaseInsensitive bool
}

func (q Branch) Match(node ast.Node, index int, source []byte) bool {
	heading, ok := node.(*preprocess.TreeBranch)
	if !ok {
		return false
	}
	if q.Level != 0 && heading.Level != q.Level {
		return false
	}

	if q.CaseInsensitive {
		return bytes.EqualFold(node.FirstChild().Text(source), q.Name)
	}
	return bytes.Equal(node.FirstChild().Text(source), q.Name)
}

// Matches an ordered or unordered list.
type List struct {
	CaseInsensitive bool
}

func (q List) Match(node ast.Node, index int, source []byte) bool {
	_, ok := node.(*ast.List)
	return ok
}

// Wraps another query, only when it's the nth child of the parent.
// Note that for branches, the heading itself is the 0th child.
type Index struct {
	Index int
	Node  Node
}

func (q Index) Match(node ast.Node, index int, source []byte) bool {
	if q.Index != index {
		return false
	}

	return q.Node.Match(node, index, source)
}
