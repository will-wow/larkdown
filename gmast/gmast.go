// gmast provides some helper functions for working with goldmark's AST.
// Useful for writing NodeUnmarshalers.
package gmast

import "github.com/yuin/goldmark/ast"

// Run a callback on each list item in a list.
func ForEachListItem(node *ast.List, source []byte, fn func(item ast.Node, index int)) {
	ForEachChild(node, source, func(child ast.Node, index int) {
		if _, ok := child.(*ast.ListItem); ok {
			fn(child, index)
		}
	})
}

// Run a callback on each direct child of a node.
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

// Get the next sibling, or the next ancestor's sibling.
func GetNextSibling(node ast.Node) ast.Node {
	if node == nil {
		return nil
	}

	next := node.NextSibling()
	if next != nil {
		return next
	}

	return getNextParentSiblingToProcess(node.Parent())
}

// Walk up the tree until we find a parent with a sibling to process.
func getNextParentSiblingToProcess(node ast.Node) ast.Node {
	if node == nil {
		return nil
	}

	parent := node.Parent()
	for {
		if parent == nil {
			return nil
		}

		next := parent.NextSibling()
		if next != nil {
			return next
		}

		parent = parent.Parent()
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
