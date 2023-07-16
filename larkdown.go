package larkdown

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/preprocess"
	"github.com/will-wow/larkdown/query"
)

// A markdown document parsed into a tree of goldmark AST nodes.
type Tree struct {
	Doc    preprocess.TreeBranch
	Source []byte
}

// Parse a markdown document into a tree of goldmark AST nodes for querying and unmarshaling.
func MarkdownToTree(source []byte) (Tree, error) {
	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(source))

	return GoldmarkToTree(doc, source)
}

// Parse an already parsed markdown document into a tree of goldmark AST nodes for querying and unmarshaling.
// Use this to apply goldmark extensions before processing.
func GoldmarkToTree(doc ast.Node, source []byte) (Tree, error) {
	tree, err := preprocess.GoldmarkToTree(doc, source)
	if err != nil {
		return Tree{}, err
	}

	return Tree{Doc: tree, Source: source}, nil
}

// Use a matcher to find a node, and then unmarshal its contents into structured data.
func Unmarshal(tree Tree, matcher []match.Node, parser NodeUnmarshaler) error {
	found, err := query.FindMatch(tree.Doc, matcher, tree.Source)
	if err != nil {
		return err
	}

	return parser.UnmarshalNode(found, tree.Source)
}
