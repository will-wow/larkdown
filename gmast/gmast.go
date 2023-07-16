// Package gmast provides some helper functions for working with goldmark's AST.
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

// Run a callback on each child of a node.
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
