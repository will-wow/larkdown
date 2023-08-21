// gmast provides some helper functions for working with goldmark's AST.
// Useful for writing NodeUnmarshalers.
package gmast

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
)

// ForEachListItem runs a callback on each list item in a list.
func ForEachListItem(node *ast.List, source []byte, fn func(item ast.Node, index int)) {
	ForEachChild(node, source, func(child ast.Node, index int) {
		if _, ok := child.(*ast.ListItem); ok {
			fn(child, index)
		}
	})
}

// ForEachChild runs a callback on each direct child of a node.
func ForEachChild(node ast.Node, source []byte, fn func(child ast.Node, index int)) {
	child := node.FirstChild()
	index := 0
	for {
		if child == nil {
			break
		}

		fn(child, index)
		index++
		child = child.NextSibling()
	}
}

// Do a depth-first walk of the AST, calling the walker function on each node, going to siblings, until the walker returns WalkStop or an error or hits EOF.
func WalkSiblingsUntil(node ast.Node, walker ast.Walker) error {
	for {
		if node == nil {
			return nil
		}

		status, err := walkHelper(node, walker)
		if err != nil || status == ast.WalkStop {
			return err
		}

		node = node.NextSibling()
	}
}

// Copied from goldmark's ast.walkHelper
func walkHelper(n ast.Node, walker ast.Walker) (ast.WalkStatus, error) {
	status, err := walker(n, true)
	if err != nil || status == ast.WalkStop {
		return status, err
	}
	if status != ast.WalkSkipChildren {
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if st, err := walkHelper(c, walker); err != nil || st == ast.WalkStop {
				return ast.WalkStop, err
			}
		}
	}
	status, err = walker(n, false)
	if err != nil || status == ast.WalkStop {
		return ast.WalkStop, err
	}
	return ast.WalkContinue, nil
}

// FindSibling finds the first direct sibling of a node that matches the given predicate.
// If no match is found, returns nil.
func FindSibling(node ast.Node, isMatch func(node ast.Node) bool) ast.Node {
	for c := node.NextSibling(); c != nil; c = c.NextSibling() {
		if isMatch(c) {
			return c
		}
	}
	return nil
}

// LastChildOfHeading returns the last node that is before
// a heading of a lower level or the end of the document.
func LastChildOfHeading(node ast.Node) (sibling ast.Node, err error) {
	heading, ok := node.(*ast.Heading)
	if !ok {
		return nil, fmt.Errorf("not a heading")
	}

	sibling = node

	var peek ast.Node

	for {
		peek = sibling.NextSibling()
		// If there is no next sibling, the current sibling is the last sibling.
		if peek == nil {
			return sibling, nil
		}
		// If the next sibling is a heading of the same or higher level,
		// the current sibling is the last sibling.
		if IsHeadingLevelBelow(peek, heading.Level) {
			return sibling, nil
		}

		// Go to the next sibling.
		sibling = peek
	}
}

// IsHeadingLevelBelow checks if a node is a heading that is outside a given level's children.
func IsHeadingLevelBelow(node ast.Node, level int) bool {
	if node.Kind() != ast.KindHeading {
		return false
	}

	heading, ok := node.(*ast.Heading)
	if !ok {
		return false
	}

	return heading.Level <= level
}
