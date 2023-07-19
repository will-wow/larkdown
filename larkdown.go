package larkdown

import (
	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
)

// // Parse a markdown document into a tree of goldmark AST nodes for querying and unmarshaling.
// func MarkdownToTree(source []byte) (preprocess.Tree, error) {
// 	md := goldmark.New()
// 	doc := md.Parser().Parse(text.NewReader(source))

// 	return GoldmarkToTree(doc, source)
// }

// // Parse an already parsed markdown document into a tree of goldmark AST nodes for querying and unmarshaling.
// // Use this to apply goldmark extensions before processing.
// func GoldmarkToTree(doc ast.Node, source []byte) (preprocess.Tree, error) {
// 	return preprocess.GoldmarkToTree(doc, source)
// }

// Use a matcher to find a node, and then unmarshal its contents into structured data.
func Unmarshal(doc ast.Node, source []byte, matcher []match.Node, parser NodeUnmarshaler) error {
	found, err := query.QueryTree(doc, source, matcher)
	if err != nil {
		return err
	}

	return parser.UnmarshalNode(found, source)
}
