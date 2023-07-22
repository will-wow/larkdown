package larkdown_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/match"
)

func TestFind(t *testing.T) {
	t.Run("Recipe", func(t *testing.T) {
		source, err := os.ReadFile("examples/recipe.md")
		require.NoError(t, err)

		md := goldmark.New()
		doc := md.Parser().Parse(text.NewReader(source))

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
		source, err := os.ReadFile("examples/simple.md")
		require.NoError(t, err)

		md := goldmark.New(
			goldmark.WithExtensions(
				&hashtag.Extender{
					// Resolver: hashtagResolver,
					Variant: hashtag.ObsidianVariant,
				},
			),
		)
		doc := md.Parser().Parse(text.NewReader(source))

		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
			match.Tag{},
		}

		tags, err := larkdown.FindAll(doc, source, matcher, func(node ast.Node, source []byte) (string, error) {
			return string(node.Text(source)), nil
		})
		require.NoError(t, err)

		require.Equal(t, []string{"#tag1", "#tag2", "#tag3"}, tags)
	})

	t.Run("Tags", func(t *testing.T) {
		source, err := os.ReadFile("examples/simple.md")
		require.NoError(t, err)

		md := goldmark.New(
			goldmark.WithExtensions(
				&hashtag.Extender{
					// Resolver: hashtagResolver,
					Variant: hashtag.ObsidianVariant,
				},
			),
		)
		doc := md.Parser().Parse(text.NewReader(source))

		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
			match.Tag{},
		}

		tags, err := larkdown.FindAll(doc, source, matcher, larkdown.DecodeTag)
		require.NoError(t, err)

		require.Equal(t, []string{"tag1", "tag2", "tag3"}, tags)
	})
}
