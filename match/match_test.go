package match_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/internal/test"
	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
)

func TestNodeOfKind(t *testing.T) {
	t.Run("Link", func(t *testing.T) {
		tree, source := test.TreeFromMd(t, `
			[link](https://example.com)
		`)

		matcher := []match.Node{
			match.NodeOfKind{Kind: ast.KindParagraph},
			match.NodeOfKind{Kind: ast.KindLink},
		}

		match, err := query.QueryOne(tree, source, matcher)
		require.NoError(t, err)
		require.Equal(t, match.Kind(), ast.KindLink)
	})
}

func TestHeading(t *testing.T) {
	t.Run("should match by level", func(t *testing.T) {
		tree, source := test.TreeFromMd(t, `
		# Heading 1

		Body

		## Heading 2

		Body2
		`)

		matcher := []match.Node{
			match.Heading{Level: 2},
		}

		match, err := query.QueryOne(tree, source, matcher)
		require.NoError(t, err)
		require.Equal(t, "Heading 2", string(match.Text(source)))
	})
}

func TestTable(t *testing.T) {
	t.Run("should match and decode table cell contents", func(t *testing.T) {
		tree, source := test.TreeFromMd(t, `
			# Table

			| number | text |
			| ------ | ---- |
			| 1      | one  |
			| 2      | two  |
			`,
			goldmark.WithExtensions(extension.Table),
		)

		matcher := []match.Node{
			match.Table{},
		}

		match, err := query.QueryOne(tree, source, matcher)
		require.NoError(t, err, "failed to find table")

		table, err := larkdown.DecodeTableToMap(match, source)
		require.NoError(t, err, "table failed to decode")

		require.Len(t, table, 2, "table should have 2 rows")
		require.Equal(
			t,
			map[string]string{"number": "1", "text": "one"},
			table[0],
			"row one should have values",
		)
		require.Equal(
			t,
			map[string]string{"number": "2", "text": "two"},
			table[1],
			"row two should have values",
		)
	})

	t.Run("should decode table with blank headers", func(t *testing.T) {
		tree, source := test.TreeFromMd(t, `
			# Table

			|        |      |
			| ------ | ---- |
			| 1      | one  |
			| 2      | two  |
			`,
			goldmark.WithExtensions(extension.Table),
		)

		matcher := []match.Node{
			match.Table{},
		}

		match, err := query.QueryOne(tree, source, matcher)
		require.NoError(t, err, "failed to find table")

		table, err := larkdown.DecodeTableToMap(match, source)
		require.NoError(t, err, "table failed to decode")

		require.Len(t, table, 2, "table should have 2 rows")
		require.Equal(
			t,
			map[string]string{"0": "1", "1": "one"},
			table[0],
			"row one should have values",
		)
		require.Equal(
			t,
			map[string]string{"0": "2", "1": "two"},
			table[1],
			"row two should have values",
		)
	})
}
