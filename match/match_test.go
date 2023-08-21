package match_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/internal/test"
	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
)

func TestNodeOfKind(t *testing.T) {
	t.Run("Link", func(t *testing.T) {
		tree, source := test.TreeFromMd(t, `
			[link](https://example.com)
		`)

		matcher := []match.Node{
			match.NodeOfKind{Kind: ast.KindParagraph},
			match.NodeOfKind{Kind: ast.KindLink},
		}

		match, err := query.QueryOne(tree, source, matcher)
		require.NoError(t, err)
		require.Equal(t, match.Kind(), ast.KindLink)
	})
}

func TestHeading(t *testing.T) {
	t.Run("should match by level", func(t *testing.T) {
		tree, source := test.TreeFromMd(t, `
		# Heading 1

		Body

		## Heading 2

		Body2
		`)

		matcher := []match.Node{
			match.Heading{Level: 2},
		}

		match, err := query.QueryOne(tree, source, matcher)
		require.NoError(t, err)
		require.Equal(t, "Heading 2", string(match.Text(source)))
	})
}
