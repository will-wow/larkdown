package query

import (
	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/gmast"
	"github.com/will-wow/larkdown/match"
)

// Use all but the last matcher in the query to find a parent node,
// then return a list of descendants that match the last matcher.
func QueryAll(doc ast.Node, source []byte, query []match.Node, extractor match.Node) (found []ast.Node, err error) {
	var node ast.Node
	var lastMatcher match.Node
	if len(query) == 0 {
		// If there are no matchers, we're matching the whole document.
		node = doc
		// And the last matcher is a dummy matcher that would have matched anything.
		lastMatcher = match.AnyNode{}
	} else {
		node, err = QueryOne(doc, source, query)
		if err != nil {
			return found, err
		}
		lastMatcher = query[len(query)-1]
	}

	return queryDescendants(node, source, extractor, lastMatcher)
}

func queryDescendants(
	matchedNode ast.Node,
	source []byte,
	extractor match.Node,
	lastMatch match.Node,
) (out []ast.Node, err error) {
	startNode := lastMatch.NextNode(matchedNode)

	err = gmast.WalkSiblingsUntil(startNode, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// Handle heading ends
		if lastMatch != nil {
			if lastMatch.EndMatch(node) {
				return ast.WalkStop, nil
			}
		}

		if extractor.Match(node, 0, source) {
			out = append(out, node)
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil

	})
	if err != nil {
		return out, err
	}
	return out, nil
}
