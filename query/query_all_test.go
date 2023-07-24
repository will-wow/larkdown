package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown/internal/test"
	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
)

func TestQueryAll(t *testing.T) {
	t.Run("QueryAll", func(t *testing.T) {
		tree, source := test.TreeFromFile(t, "../examples/simple.md",
			goldmark.WithExtensions(
				&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			),
		)

		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
		}

		matches, err := query.QueryAll(tree, source, matcher, match.Tag{})
		require.NoError(t, err)
		require.Equal(t, 3, len(matches))
	})

	t.Run("no parent", func(t *testing.T) {
		tree, source := test.TreeFromFile(t, "../examples/simple.md",
			goldmark.WithExtensions(
				&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			),
		)

		matcher := []match.Node{}

		matches, err := query.QueryAll(tree, source, matcher, match.Tag{})
		require.NoError(t, err)
		require.Equal(t, 3, len(matches))
	})
}
