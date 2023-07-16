package parser

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"
)

// TODO: queries:
// Code block with language
// Nth Instance Query that retunrns a "don't pop this query", and keeps state of how many times it's called?

type NodeQuery interface {
	Match(node ast.Node, index int, source []byte) (ok bool)
}

// Matches a heading by level and name.
type BranchQuery struct {
	Level           int
	Name            []byte
	CaseInsensitive bool
}

func (q BranchQuery) Match(node ast.Node, index int, source []byte) bool {
	heading, ok := node.(*TreeBranch)
	if !ok {
		return false
	}
	if q.Level != 0 && heading.Level != q.Level {
		return false
	}

	if q.CaseInsensitive {
		return bytes.EqualFold(node.FirstChild().Text(source), q.Name)
	}
	return bytes.Equal(node.FirstChild().Text(source), q.Name)
}

// Matches an ordered or unordered list.
type ListQuery struct {
	CaseInsensitive bool
}

func (q ListQuery) Match(node ast.Node, index int, source []byte) bool {
	_, ok := node.(*ast.List)
	return ok
}

// Wraps another query, only when it's the nth child of the parent.
// Note that for branches, the heading itself is the 0th child.
type IndexQuery struct {
	Index int
	Query NodeQuery
}

func (q IndexQuery) Match(node ast.Node, index int, source []byte) bool {
	if q.Index != index {
		return false
	}

	return q.Query.Match(node, index, source)
}

type ResultParser interface {
	ParseResult(node ast.Node, source []byte) error
}

func ParseResult(match ast.Node, source []byte, parser ResultParser) error {
	return parser.ParseResult(match, source)
}

type ListParser struct {
	Items []string
}

func (p *ListParser) ParseResult(node ast.Node, source []byte) error {
	list, ok := node.(*ast.List)
	if !ok {
		return fmt.Errorf("expected list node")
	}

	ForEachListItem(list, source, func(item ast.Node, _ int) {
		p.Items = append(p.Items, string(item.Text(source)))
	})

	return nil
}

func ForEachListItem(node *ast.List, source []byte, fn func(item ast.Node, index int)) {
	forEachChild(node, source, func(child ast.Node, index int) {
		if _, ok := child.(*ast.ListItem); ok {
			fn(child, index)
		}
	})

}

func forEachChild(node ast.Node, source []byte, fn func(child ast.Node, index int)) {
	child := node.FirstChild()
	index := 0
	for {
		if child == nil {
			break
		}

		fn(child, index)
		index++
		child = child.NextSibling()
	}
}

func FindMatch(doc TreeBranch, query []NodeQuery, source []byte) (ast.Node, error) {
	queryCount := len(query)

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
			fmt.Println("eof")
			break
		}

		fmt.Printf("node: %s\n", string(node.Text(source)))

		// If we are out of queries, failure. Break.
		if activeQueryIndex == queryCount {
			break
		}

		match := query[activeQueryIndex].Match(node, queryChildIndex, source)
		if !match {
			fmt.Println("no match", activeQueryIndex)
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

// Parse parses a markdown document into a Heading tree
func MarkdownToTree(source []byte) (TreeBranch, error) {
	md := goldmark.New(
		goldmark.WithExtensions(&hashtag.Extender{}),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	doc := md.Parser().Parse(text.NewReader(source))

	tree := NewTreeBranchRoot()

	// The active branch, which can move down as we go through levels
	activeTreeBranch := tree

	next := doc.FirstChild()
	for {
		child := next
		// If we are at the end of the document, break out of the loop
		if child == nil {
			break
		}

		// Record the next sibling, becase the child is about to be moved.
		next = child.NextSibling()

		switch node := child.(type) {
		case *ast.Heading:
			// Go up to the first parent of this heading that is at a lower level
			// This will not change the active branch if the heading is at a lower level
			activeTreeBranch = findParentBeforeLevel(activeTreeBranch, node.Level)

			// Create a new TreeHeading,
			// with the parent being the active heading
			// and the first child being the real heading
			treeHeading := NewTreeBranch(node, activeTreeBranch)

			appendToBranch(activeTreeBranch, treeHeading)

			// Note that the new heading is now the active heading
			activeTreeBranch = treeHeading
		default:
			appendToBranch(activeTreeBranch, node)
		}
	}

	return *tree, nil
}

// If we are not at the root of the tree,
// append the new level to the active level
func findParentBeforeLevel(activeTreeLevel *TreeBranch, level int) *TreeBranch {
	if level < 1 {
		panic("level must be greater than 0")
	}
	for {
		if activeTreeLevel == nil {
			panic("missed the root of the tree")
		}

		if activeTreeLevel.Level < level {
			return activeTreeLevel
		}

		activeTreeLevel = activeTreeLevel.TreeParent
	}
}

// If we are not at the root of the tree,
// append the new level to the active level
func appendToBranch(activeTreeBranch *TreeBranch, node ast.Node) {
	activeTreeBranch.AppendChild(activeTreeBranch, node)
}

type TreeBranch struct {
	ast.BaseInline
	TreeParent *TreeBranch
	Level      int
}

func (n *TreeBranch) Dump(source []byte, level int) {
}

var KindTreeBranch = ast.NewNodeKind("TreeBranch")

func (n *TreeBranch) Kind() ast.NodeKind {
	return KindTreeBranch
}

// A tree branch for the root of the document.
func NewTreeBranchRoot() *TreeBranch {
	return &TreeBranch{
		TreeParent: nil,
		Level:      0,
		BaseInline: ast.BaseInline{},
	}
}

func NewTreeBranch(heading *ast.Heading, parent *TreeBranch) *TreeBranch {
	if heading == nil {
		panic("heading cannot be nil")
	}

	headingContents := *heading
	threeBranch := &TreeBranch{
		TreeParent: parent,
		Level:      headingContents.Level,
		BaseInline: ast.BaseInline{},
	}

	threeBranch.AppendChild(threeBranch, heading)

	return threeBranch
}
