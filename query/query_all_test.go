package query_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
)

func TestQueryAll(t *testing.T) {
	t.Run("QueryAll", func(t *testing.T) {
		source, err := os.ReadFile("../examples/simple.md")
		require.NoError(t, err)

		md := goldmark.New(
			goldmark.WithExtensions(
				&hashtag.Extender{
					// Resolver: hashtagResolver,
					Variant: hashtag.ObsidianVariant,
				},
			),
		)
		tree := md.Parser().Parse(text.NewReader(source))

		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
			match.Tag{},
		}

		matches, err := query.QueryAll(tree, source, matcher)
		require.NoError(t, err)
		require.Equal(t, 3, len(matches))
	})

	t.Run("no parent", func(t *testing.T) {
		source, err := os.ReadFile("../examples/simple.md")
		require.NoError(t, err)

		md := goldmark.New(
			goldmark.WithExtensions(
				&hashtag.Extender{
					// Resolver: hashtagResolver,
					Variant: hashtag.ObsidianVariant,
				},
			),
		)
		tree := md.Parser().Parse(text.NewReader(source))

		matcher := []match.Node{
			match.Tag{},
		}

		matches, err := query.QueryAll(tree, source, matcher)
		require.NoError(t, err)
		require.Equal(t, 3, len(matches))
	})
}
