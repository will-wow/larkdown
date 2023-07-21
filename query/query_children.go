package query

import (
	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/match"
)

func QueryChildren(parentNode ast.Node, source []byte, matcher match.SearchFor) (out []ast.Node, err error) {
	err = ast.Walk(parentNode, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if matcher.Node.Match(node, 0, source) {
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
