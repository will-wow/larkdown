package examples_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/will-wow/larkdown/examples"
)

func TestReadme(t *testing.T) {
	t.Run("ParseRecipe works", func(t *testing.T) {
		recipe, err := examples.ParseRecipe()
		require.NoError(t, err)

		require.Equal(t, []string{"Chicken", "Vegetables", "Salt", "Pepper"}, recipe.Ingredients)
		require.Equal(t, []string{"dinner", "chicken"}, recipe.Tags)
		require.Contains(t, recipe.Html.String(), "Here's a long story about making dinner.")
	})
}
