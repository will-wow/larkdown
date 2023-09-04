package gmast_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/gmast"
	"github.com/will-wow/larkdown/internal/test"
)

func TestAppendChild(t *testing.T) {
	source := []byte("")

	newSegment, source := gmast.NewSegment("Hello, world!", source)

	list := gmast.AppendChild(ast.NewList('-'),
		gmast.AppendChild(ast.NewListItem(2),
			gmast.AppendChild(
				ast.NewTextBlock(),
				ast.NewTextSegment(newSegment),
			)))

	// Added the new text to the source
	require.Contains(t, string(source), "Hello, world!")
	require.Equal(t, list.Kind(), ast.KindList)
}

func TestAppendHeadingChild(t *testing.T) {
	tree, source := test.TreeFromMd(t, `
	## First Heading

	P1

	P2

	## Second Heading

	P3
	`)

	// Append P3 to the end of ## First Heading
	newParagraph, source := gmast.NewParagraph("P3", source)
	gmast.AppendHeadingChild(tree.FirstChild(), newParagraph)

	// Render back to markdown
	var renderedMd bytes.Buffer
	err := larkdown.NewNodeRenderer().Render(&renderedMd, source, tree)
	require.NoError(t, err, "error rendering back to markdown")

	require.Contains(t, renderedMd.String(), "P2\n\nP3", "P3 is not after P2 in MD")

	var renderedHtml bytes.Buffer
	err = goldmark.DefaultRenderer().Render(&renderedHtml, source, tree)
	require.NoError(t, err, "error rendering to HTML")

	require.Contains(t, renderedHtml.String(), "<p>P2</p>\n<p>P3</p>", "P3 is not after P2 in HTML")
}
