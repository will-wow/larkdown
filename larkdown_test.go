package larkdown_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/internal/test"
	"github.com/will-wow/larkdown/match"
)

func TestFind(t *testing.T) {
	t.Run("Recipe", func(t *testing.T) {
		doc, source := test.TreeFromFile(t, "examples/recipe.md")

		matcher := []match.Node{
			match.Branch{Level: 2, Name: []byte("Ingredients")},
			match.Branch{Level: 3, Name: []byte("Buy")},
			match.Index{Index: 0, Node: match.List{}},
		}

		list, err := larkdown.Find(doc, source, matcher, larkdown.DecodeListItems)
		require.NoError(t, err)

		require.Equal(t, []string{"1 Medium Apple", "1 small-medium carrot", "1 banana", "2 eggs"}, list)
	})
}

func TestFindAll(t *testing.T) {
	t.Run("Custom decoder", func(t *testing.T) {
		doc, source := test.TreeFromFile(t, "examples/simple.md",
			goldmark.WithExtensions(
				&hashtag.Extender{
					Variant: hashtag.ObsidianVariant,
				},
			),
		)

		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
		}

		tags, err := larkdown.FindAll(
			doc, source, matcher,
			match.Tag{},
			func(node ast.Node, source []byte) (string, error) {
				return string(node.Text(source)), nil
			})
		require.NoError(t, err)

		require.Equal(t, []string{"#tag1", "#tag2", "#tag3"}, tags)
	})

	t.Run("Tags", func(t *testing.T) {
		doc, source := test.TreeFromFile(t, "examples/simple.md",
			goldmark.WithExtensions(
				&hashtag.Extender{
					Variant: hashtag.ObsidianVariant,
				},
			),
		)

		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
		}

		tags, err := larkdown.FindAll(doc, source, matcher, match.Tag{}, larkdown.DecodeTag)
		require.NoError(t, err)

		require.Equal(t, []string{"tag1", "tag2", "tag3"}, tags)
	})
}

func BenchmarkFind(b *testing.B) {
	doc, source := test.TreeFromFile(b, "examples/recipe.md")

	matcher := []match.Node{
		match.Branch{Level: 2, Name: []byte("Ingredients")},
		match.Branch{Level: 3, Name: []byte("Buy")},
		match.Index{Index: 0, Node: match.List{}},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = larkdown.Find(doc, source, matcher, larkdown.DecodeListItems)
	}
}

func BenchmarkFindAll(b *testing.B) {
	doc, source := test.TreeFromFile(b, "examples/simple.md",
		goldmark.WithExtensions(
			&hashtag.Extender{
				Variant: hashtag.ObsidianVariant,
			},
		),
	)

	matcher := []match.Node{
		match.Branch{Level: 1, Name: []byte("Title")},
		match.Branch{Level: 2, Name: []byte("Subheading")},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = larkdown.FindAll(doc, source, matcher, match.Tag{}, larkdown.DecodeTag)
	}
}
