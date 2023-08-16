package markdown_test

import (
	"bytes"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/will-wow/larkdown/renderer/markdown"
)

func NewRenderer() renderer.Renderer {
	return renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(markdown.NewRenderer(), 1000)))
}

func TestSimpleRenderer(t *testing.T) {
	contents := `
	# Heading 1

	Paragraph is *here*
	and **here** as well

	- first item
	- second item

	1. First number
	1. Second number

	---

	[a link](http://example.com)

	<http://example.com>
	`

	source := []byte(dedent.Dedent(contents)[1:])

	md := goldmark.New(
		goldmark.WithRenderer(NewRenderer()),
	)
	doc := md.Parser().Parse(text.NewReader(source))

	var rendered bytes.Buffer

	err := md.Renderer().Render(&rendered, source, doc)
	require.NoError(t, err)

	// doc.Dump(source, 3)

	require.Equal(t, string(source), rendered.String())
}
