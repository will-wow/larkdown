package larkdown

import (
	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
)

// Use a matcher to find a node, and then unmarshal its contents into structured data.
func Unmarshal(doc ast.Node, source []byte, matcher []match.Node, parser NodeUnmarshaler) error {
	found, err := query.QueryTree(doc, source, matcher)
	if err != nil {
		return err
	}

	return parser.UnmarshalNode(found, source)
}
