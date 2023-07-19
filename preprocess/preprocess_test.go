package preprocess_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"

	"github.com/will-wow/larkdown/preprocess"
)

func TestGoldmarkToTree(t *testing.T) {
	t.Run("should not change the output of goldmark rendering", func(t *testing.T) {
		source := []byte(md)

		md := goldmark.New()
		doc := md.Parser().Parse(text.NewReader(source))

		var b bytes.Buffer
		err := md.Renderer().Render(&b, source, doc)
		require.NoError(t, err)
		require.Equal(t, string(want), b.String())

		_, _ = preprocess.GoldmarkToTree(doc, source)

		var b2 bytes.Buffer
		err = md.Renderer().Render(&b2, source, doc)
		require.NoError(t, err)

		require.Equal(t, string(want), b2.String())
	})
}

var md = `
# Heading 1

Hello

## Heading 2

World
`

var want = `<h1>Heading 1</h1>
<p>Hello</p>
<h2>Heading 2</h2>
<p>World</p>
`
