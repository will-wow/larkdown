package mdrender_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/gmast"
	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/mdfront"
	"github.com/will-wow/larkdown/mdrender"
	"github.com/will-wow/larkdown/query"
)

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

func TestStandaloneRenderer(t *testing.T) {
	source, err := os.ReadFile("../examples/all-tags.md")
	require.NoError(t, err, "error reading markdown file")

	md := goldmark.New(
		goldmark.WithExtensions(
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
		),
	)

	doc := md.Parser().Parse(text.NewReader(source))

	var rendered bytes.Buffer
	err = larkdown.NewNodeRenderer().Render(&rendered, source, doc)

	require.NoError(t, err)

	// Print the ast if the test is going to fail
	if string(source) != rendered.String() {
		doc.Dump(source, 4)
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

	newSegment, source := gmast.NewSegment("Hello, world!", source)

	// Add a new list item to the first list
	gmast.AppendChild(listNode,
		gmast.AppendChild(ast.NewListItem(2),
			gmast.AppendChild(
				ast.NewTextBlock(),
				ast.NewTextSegment(newSegment),
			)))

	// Render
	var rendered bytes.Buffer
	err = md.Renderer().Render(&rendered, source, doc)
	require.NoError(t, err)

	// Print the ast for debugging
	// doc.Dump(source, 3)

	// Check that the new text is in the rendered markdown at the right position
	require.Contains(t, rendered.String(), "- first item\n- second item\n- Hello, world!\n")
}

func TestFrontMatter(t *testing.T) {
	sourceMd := `
	---
	title: My Recipe
	---

	# Title

	`

	source := []byte(strings.TrimLeft(dedent.Dedent(sourceMd), "\n"))

	type metaData struct {
		Title string `yaml:"title"`
	}

	data := &metaData{}

	md := goldmark.New(
		goldmark.WithExtensions(
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			// Frontmatter parsing.
			&frontmatter.Extender{},
			// Frontmatter rendering.
			&mdfront.Extender{},
		),
		// Store a pointer to the struct now for later rendering.
		goldmark.WithRenderer(larkdown.NewNodeRenderer(mdrender.WithFrontmatter(data))),
	)

	ctx := parser.NewContext()
	doc := md.Parser().Parse(text.NewReader(source), parser.WithContext(ctx))

	// Populate the meta struct
	frontmatterData := frontmatter.Get(ctx)
	err := frontmatterData.Decode(data)
	require.NoError(t, err, "error decoding frontmatter")

	var rendered bytes.Buffer
	err = md.Renderer().Render(&rendered, source, doc)
	require.NoError(t, err)

	// Print the ast if the test is going to fail
	if string(source) != rendered.String() {
		doc.Dump(source, 3)
	}

	require.Equal(t, string(source), rendered.String())
}

func setup(t *testing.T) (source []byte, md goldmark.Markdown, doc ast.Node) {
	source, err := os.ReadFile("../examples/all-tags.md")
	require.NoError(t, err, "error reading markdown file")

	md = goldmark.New(
		goldmark.WithExtensions(
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
		),
		goldmark.WithRenderer(larkdown.NewNodeRenderer()),
	)

	doc = md.Parser().Parse(text.NewReader(source))

	return source, md, doc
}
