package query

import (
	"fmt"

	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/gmast"
	"github.com/will-wow/larkdown/match"
)

func QueryAll(doc ast.Node, source []byte, matcher []match.Node) (found []ast.Node, err error) {
	if len(matcher) == 0 {
		return found, fmt.Errorf("No matcher provided")
	}

	var node ast.Node
	var lastMatcher match.Node
	if len(matcher) == 1 {
		// If there is only one matcher, we're matching the whole document.
		node = doc
		// And the last matcher is a dummy matcher that would have matched anything.
		lastMatcher = match.AnyNode{}
	} else {
		node, err = QueryTree(doc, source, matcher[:len(matcher)-1])
		if err != nil {
			return found, err
		}
		lastMatcher = matcher[len(matcher)-2]
	}

	extractor := matcher[len(matcher)-1]

	return queryChildren(node, source, extractor, lastMatcher)
}

func queryChildren(matchedNode ast.Node, source []byte, extractor match.Node, lastMatch match.Node) (out []ast.Node, err error) {
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
