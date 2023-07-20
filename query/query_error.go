package query

import (
	"bytes"
	"fmt"

	"github.com/will-wow/larkdown/match"
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
