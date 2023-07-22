// query handles finding a match in a tree, but not unmarshaling the node.
// This package should not generally be used directly.
package query

import (
	"fmt"

	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/gmast"
	"github.com/will-wow/larkdown/match"
)

// Apply a matcher to a tree, and return the matching node for unmarshaling.
func QueryTree(doc ast.Node, source []byte, matcher []match.Node) (found ast.Node, err error) {
	queryCount := len(matcher)

	if queryCount == 0 {
		return nil, fmt.Errorf("no queries provided")
	}

	// Tracks how far we are in looping through the queries
	activeQueryIndex := 0
	// Tracks how many nodes have been processed since the last query match. This allows for index queries.
	queryChildIndex := 0
	// Tracks the active branch/heading query.
	// This is important because headings do not have children, so we need to know the active heading to know
	// when a new heading of a higher level is encountered, which stops conceptual heading block.
	var activeBranch match.Node

	// An error message that gathers all the valid matches. Only if the query fails will this be returned.
	queryError := newQueryError(queryCount)

	// Start processing at the first node of the document.
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

		if activeBranch != nil {
			if activeBranch.EndMatch(node) {
				break
			}
		}

		matcher := matcher[activeQueryIndex]

		match := matcher.Match(node, queryChildIndex, source)
		if !match {
			node = gmast.GetNextSibling(node)
			queryChildIndex++
			continue
		}

		queryError.addMatch(matcher)

		if matcher.IsFlatBranch() {
			activeBranch = matcher
		}

		// If we have a query match, then:

		// If we are not at the last query:
		if (activeQueryIndex) < queryCount-1 {
			// Either go down a level, or go to the next sibling
			node = matcher.NextNode(node)

			// go to the next query
			activeQueryIndex++
			// Reset the child index so index queries restart at 0
			queryChildIndex = 0
			continue
		}

		// Success!
		return node, nil
	}

	// Add the last failed match the error
	queryError.addFailedMatch(matcher[activeQueryIndex])

	// Return the error with the list of good matches and the bad match
	return nil, queryError
}
