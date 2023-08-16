package markdown_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
	"github.com/will-wow/larkdown/renderer/markdown"
)

func NewRenderer() renderer.Renderer {
	return renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(markdown.NewRenderer(), 998)))
}

func TestSimpleRenderer(t *testing.T) {
	source, md, doc := setup(t)

	var rendered bytes.Buffer
	err := md.Renderer().Render(&rendered, source, doc)
	require.NoError(t, err)

	// Print the ast if the test is going to fail
	if string(source) != rendered.String() {
		doc.Dump(source, 3)
	}

	require.Equal(t, string(source), rendered.String())
}

func TestAddingANode(t *testing.T) {
	source, md, doc := setup(t)

	// Add a node to the document
	firstListQuery := []match.Node{
		match.Branch{Level: 2, Name: []byte("Heading 2")},
		match.Index{Index: 0, Node: match.List{}},
	}

	listNode, err := query.QueryOne(doc, source, firstListQuery)
	require.NoError(t, err, "error finding first list")
	_, ok := listNode.(*ast.List)
	require.True(t, ok, "first list is not a list")

	newText := "Hello, world!"
	newSegment := text.NewSegment(len(source), len(source)+len(newText))
	// Add new text on the bottom of source for referencing by the segment
	source = append(source, newText...)
	newSegmentNode := ast.NewTextSegment(newSegment)

	newTextBlock := ast.NewTextBlock()
	newTextBlock.AppendChild(newTextBlock, newSegmentNode)

	newListItem := ast.NewListItem(2)
	newListItem.AppendChild(newListItem, newTextBlock)

	listNode.AppendChild(listNode, newListItem)

	// Render
	var rendered bytes.Buffer
	err = md.Renderer().Render(&rendered, source, doc)
	require.NoError(t, err)

	// Print the ast for debugging
	// doc.Dump(source, 3)

	// Check that the new text is in the rendered markdown at the right position
	require.Contains(t, rendered.String(), "- first item\n- second item\n- Hello, world!\n")
}

func setup(t *testing.T) (source []byte, md goldmark.Markdown, doc ast.Node) {
	source, err := os.ReadFile("../../examples/all-tags.md")
	require.NoError(t, err, "error reading markdown file")

	md = goldmark.New(
		goldmark.WithExtensions(
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
		),
		goldmark.WithRenderer(NewRenderer()),
	)

	doc = md.Parser().Parse(text.NewReader(source))

	return source, md, doc
}
