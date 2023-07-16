// query handles finding a match in a tree, but not unmarshaling the node.
// This package should not generally be used directly.
package query

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/preprocess"
)

// Error returned when a query fails to match.
// Includes the list of matches that were found, and the match that failed.
// Prints an error message that can be used to debug the query.
type QueryError struct {
	Matches     []match.Node
	FailedMatch match.Node
}

func newQueryError(matcherLength int) *QueryError {
	return &QueryError{
		Matches:     make([]match.Node, 0, matcherLength),
		FailedMatch: nil,
	}
}

func (e *QueryError) addMatch(node match.Node) {
	e.Matches = append(e.Matches, node)
}

func (e *QueryError) addFailedMatch(node match.Node) {
	e.FailedMatch = node
}

func (e *QueryError) Error() string {
	var matches bytes.Buffer

	matches.WriteString("document")

	for _, match := range e.Matches {
		matches.WriteString(match.String())
	}

	return fmt.Sprintf("failed to match query: %s did not have a %s", matches.String(), e.FailedMatch)
}

// Apply a matcher to a tree, and return the matching node for unmarshaling.
func QueryTree(tree preprocess.Tree, matcher []match.Node) (found ast.Node, err error) {
	doc := tree.Doc
	source := tree.Source

	queryCount := len(matcher)

	if queryCount == 0 {
		return nil, fmt.Errorf("no queries provided")
	}

	activeQueryIndex := 0
	queryChildIndex := 0

	queryError := newQueryError(queryCount)

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

		queryError.addMatch(matcher[activeQueryIndex])

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

	// Add the last failed match the error
	queryError.addFailedMatch(matcher[activeQueryIndex])

	// Return the error with the list of good matches and the bad match
	return nil, queryError
}

// func matchesForError(matches []match.Node) string {
// 	return fmt.Sprintf("matches: %+v", matches)
// }

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
