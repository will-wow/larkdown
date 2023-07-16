package larkdown_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/match"
)

func TestParse(t *testing.T) {
	t.Run("ParseRecipe", func(t *testing.T) {
		source, err := os.ReadFile("examples/recipe.md")
		require.NoError(t, err)

		tree, err := larkdown.MarkdownToTree(source)
		require.NoError(t, err)

		matcher := []match.Node{
			match.Branch{Level: 2, Name: []byte("Ingredients")},
			match.Branch{Level: 3, Name: []byte("Buy")},
			match.Index{Index: 1, Node: match.List{}},
		}

		list := &larkdown.StringList{}
		err = larkdown.Unmarshal(tree, matcher, list)
		require.NoError(t, err)

		require.Equal(t, []string{"1 Medium Apple", "1 small-medium carrot", "1 banana", "2 eggs"}, list.Items)
	})
}
