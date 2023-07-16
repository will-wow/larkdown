// Package preprocess converts a markdown document into a tree of headings
// for easy querying.
package preprocess

import (
	"github.com/yuin/goldmark/ast"
)

// Converts an already parsed goldmark ast into a Heading tree
func GoldmarkToTree(doc ast.Node, source []byte) (TreeBranch, error) {
	tree := newTreeBranchRoot()

	// The active branch, which can move down as we go through levels
	activeTreeBranch := tree

	next := doc.FirstChild()
	for {
		child := next
		// If we are at the end of the document, break out of the loop
		if child == nil {
			break
		}

		// Record the next sibling, because the child is about to be moved.
		next = child.NextSibling()

		switch node := child.(type) {
		case *ast.Heading:
			// Go up to the first parent of this heading that is at a lower level,
			// if the new heading is outside of the active branch.
			// This will not change the active branch if the current heading is at a lower level
			activeTreeBranch = findParentBeforeLevel(activeTreeBranch, node.Level)

			// Create a new TreeHeading,
			// with the parent being the active heading
			// and the first child being the real heading
			treeHeading := newTreeBranch(node, activeTreeBranch)

			// Make the new heading a child of the active heading.
			activeTreeBranch.AppendChild(activeTreeBranch, treeHeading)

			// Note that the new heading is now the active heading.
			activeTreeBranch = treeHeading
		default:
			activeTreeBranch.AppendChild(activeTreeBranch, node)
		}
	}

	return *tree, nil
}

// If we are not at the root of the tree,
// append the new level to the active level
func findParentBeforeLevel(activeTreeLevel *TreeBranch, level int) *TreeBranch {
	if level < 1 {
		panic("level must be greater than 0")
	}
	for {
		if activeTreeLevel == nil {
			panic("missed the root of the tree")
		}

		if activeTreeLevel.Level < level {
			return activeTreeLevel
		}

		activeTreeLevel = activeTreeLevel.TreeParent
	}
}
