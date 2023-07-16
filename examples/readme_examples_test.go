package examples_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/will-wow/larkdown/examples"
)

func TestReadme(t *testing.T) {
	t.Run("ParseFile works", func(t *testing.T) {
		list, err := examples.ParseFile("./recipe.md")
		require.NoError(t, err)

		require.Equal(t, []string{"1 Medium Apple", "1 small-medium carrot", "1 banana", "2 eggs"}, list)
	})
}
