// query handles finding a match in a tree, but not unmarshalling the node.
// This package should not generally be used directly.
package query

import (
	"fmt"

	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/preprocess"
)

// Apply a matcher to a tree, and return the matching node for unmarshalling.
func FindMatch(doc preprocess.TreeBranch, matcher []match.Node, source []byte) (ast.Node, error) {
	queryCount := len(matcher)

	if queryCount == 0 {
		return nil, fmt.Errorf("no queries provided")
	}

	activeQueryIndex := 0
	queryChildIndex := 0

	node := doc.FirstChild()

	if node == nil {
		return nil, fmt.Errorf("empty markdown file")
	}

	for {
		// If we are at the end of the document, failure. Break.
		if node == nil {
			break
		}

		// If we are out of queries, failure. Break.
		if activeQueryIndex == queryCount {
			break
		}

		match := matcher[activeQueryIndex].Match(node, queryChildIndex, source)
		if !match {
			node = getNextNodeToProcess(node)
			queryChildIndex++
			continue
		}

		// If we have a query match, then:

		// If we are not at the last query:
		if (activeQueryIndex) < queryCount-1 {
			// go to the next query
			activeQueryIndex++
			// Reset the child index so index queries restart at 0
			queryChildIndex = 0
			// And make the next child the first child of this element,
			node = node.FirstChild()
			continue
		}

		// Success!
		return node, nil
	}

	// TODO: Record all the query matches, so they can be used to provide context
	return nil, fmt.Errorf("no match")
}

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

func getNextNodeToProcess(node ast.Node) ast.Node {
	if node == nil {
		return nil
	}

	next := node.NextSibling()
	if next != nil {
		return next
	}

	return getNextParentSiblingToProcess(node.Parent())
}
