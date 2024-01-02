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

	t.Run("Errors when there is no match", func(t *testing.T) {
		doc, source := test.TreeFromMd(t, `# Test`)

		matcher := []match.Node{match.Heading{Level: 2}}
		_, err := larkdown.Find(doc, source, matcher, larkdown.DecodeText)
		require.ErrorContains(t, err, "failed to match query")
	})

	t.Run("returns no error when there is no match but allowNoMatch is on", func(t *testing.T) {
		doc, source := test.TreeFromMd(t, `# Test`)

		matcher := []match.Node{match.Heading{Level: 2}}
		out, err := larkdown.Find(doc, source, matcher, larkdown.DecodeText, larkdown.FindAllowNoMatch())
		require.NoError(t, err, "passes with no results")
		require.Equal(t, "", out, "returns empty value")
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
func TestFindAllErrors(t *testing.T) {
	doc, source := test.TreeFromMd(t, `# Test`, goldmark.WithExtensions(
		&hashtag.Extender{
			Variant: hashtag.ObsidianVariant,
		},
	))

	t.Run("error when there is no match from the matcher", func(t *testing.T) {
		matcher := []match.Node{match.Heading{Level: 2}}
		_, err := larkdown.FindAll(doc, source, matcher, match.Tag{}, larkdown.DecodeTag)
		require.ErrorContains(t, err, "failed to match query")
	})

	t.Run("no error when there is no match from the extractor", func(t *testing.T) {
		matcher := []match.Node{match.Heading{Level: 1}}
		out, err := larkdown.FindAll(doc, source, matcher, match.Tag{}, larkdown.DecodeTag)
		require.NoError(t, err, "failed to match query")
		require.Equal(t, []string{}, out)
	})

	t.Run("no error when there is no match from the matcher but AllowNoMatch is on", func(t *testing.T) {
		matcher := []match.Node{match.Heading{Level: 2}}
		out, err := larkdown.FindAll(doc, source, matcher, match.Tag{}, larkdown.DecodeTag, larkdown.FindAllAllowNoMatch())
		require.NoError(t, err, "failed to match query")
		require.Nil(t, out, "returns nil")
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
