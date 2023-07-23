package examples_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/will-wow/larkdown/examples"
)

func TestReadme(t *testing.T) {
	t.Run("ParseFile works", func(t *testing.T) {
		list, err := examples.ParseFile()
		require.NoError(t, err)

		require.Equal(t, []string{"Chicken", "Vegetables", "Salt", "Pepper"}, list)
	})
}
