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

func TestQueryOneErrorMessage(t *testing.T) {
	tree, source := test.TreeFromFile(t, "../examples/simple.md",
		goldmark.WithExtensions(
			&hashtag.Extender{
				Variant: hashtag.ObsidianVariant,
			},
		),
	)

	t.Run("no match", func(t *testing.T) {
		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
			// No match
			match.Branch{Level: 3, Name: []byte("Not Real")},
		}
		_, err := query.QueryOne(tree, source, matcher)
		require.ErrorContains(t, err, "failed to match query: document[# Title][## Subheading] did not have a [### Not Real]")
	})

	t.Run("wrong level", func(t *testing.T) {
		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			// Level is wrong
			match.Branch{Level: 3, Name: []byte("Subheading")},
		}
		_, err := query.QueryOne(tree, source, matcher)
		require.ErrorContains(t, err, "failed to match query: document[# Title] did not have a [### Subheading]")
	})

	t.Run("no error on missed heading", func(t *testing.T) {
		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			// There's a heading 2 in between in the MD file
			match.Branch{Level: 3, Name: []byte("Sub-subheading")},
		}
		_, err := query.QueryOne(tree, source, matcher)
		require.NoError(t, err)
	})

	t.Run("extra list", func(t *testing.T) {
		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
			match.List{},
			// Extra list
			match.List{},
		}
		_, err := query.QueryOne(tree, source, matcher)
		require.ErrorContains(t, err, "failed to match query: document[# Title][## Subheading].list did not have a .list")

	})

	t.Run("bad index", func(t *testing.T) {
		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
			match.List{},
			// Bad index
			match.Index{Index: 4, Node: match.AnyNode{}},
		}
		_, err := query.QueryOne(tree, source, matcher)
		require.ErrorContains(t, err, "failed to match query: document[# Title][## Subheading].list did not have a [4].any")
	})
}
