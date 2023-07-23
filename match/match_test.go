package match_test

import (
	"fmt"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
)

func TestNodeOfKind(t *testing.T) {
	t.Run("Heading", func(t *testing.T) {
		mdSource := `
		[link](https://example.com)
		`
		source := []byte(dedent.Dedent(mdSource))

		fmt.Println(dedent.Dedent(mdSource))

		md := goldmark.New()
		tree := md.Parser().Parse(text.NewReader(source))

		matcher := []match.Node{
			match.NodeOfKind{Kind: ast.KindParagraph},
			match.NodeOfKind{Kind: ast.KindLink},
		}

		match, err := query.QueryTree(tree, source, matcher)
		require.NoError(t, err)
		require.Equal(t, match.Kind(), ast.KindLink)
	})
}
