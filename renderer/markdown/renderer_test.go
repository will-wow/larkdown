package markdown_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown/renderer/markdown"
)

func NewRenderer() renderer.Renderer {
	return renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(markdown.NewRenderer(), 998)))
}

func TestSimpleRenderer(t *testing.T) {
	source, err := os.ReadFile("../../examples/all-tags.md")
	require.NoError(t, err, "error reading markdown file")

	md := goldmark.New(
		goldmark.WithExtensions(
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
		),
		goldmark.WithRenderer(NewRenderer()),
	)
	doc := md.Parser().Parse(text.NewReader(source))

	var rendered bytes.Buffer

	err = md.Renderer().Render(&rendered, source, doc)
	require.NoError(t, err)

	// Print the ast if the test is going to fail
	if string(source) != rendered.String() {
		doc.Dump(source, 3)
	}

	require.Equal(t, string(source), rendered.String())
}
