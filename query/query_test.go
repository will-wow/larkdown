package query_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
)

func TestQueryTree(t *testing.T) {
	t.Run("Match error messages", func(t *testing.T) {

		source, err := os.ReadFile("../examples/simple.md")
		require.NoError(t, err)

		tree, err := larkdown.MarkdownToTree(source)
		require.NoError(t, err)

		matcher := []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
			// No match
			match.Branch{Level: 3, Name: []byte("Not Real")},
		}
		_, err = query.QueryTree(tree, matcher)
		require.ErrorContains(t, err, "failed to match query: document[# Title][## Subheading] did not have a [### Not Real]")

		matcher = []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			// Level is wrong
			match.Branch{Level: 3, Name: []byte("Subheading")},
		}
		_, err = query.QueryTree(tree, matcher)
		require.ErrorContains(t, err, "failed to match query: document[# Title] did not have a [### Subheading]")

		matcher = []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			// Missed a sublevel
			match.Branch{Level: 3, Name: []byte("Sub-subheading")},
		}
		_, err = query.QueryTree(tree, matcher)
		require.ErrorContains(t, err, "failed to match query: document[# Title] did not have a [### Sub-subheading]")

		matcher = []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
			match.List{},
			// Extra list
			match.List{},
		}
		_, err = query.QueryTree(tree, matcher)
		require.ErrorContains(t, err, "failed to match query: document[# Title][## Subheading].list did not have a .list")

		matcher = []match.Node{
			match.Branch{Level: 1, Name: []byte("Title")},
			match.Branch{Level: 2, Name: []byte("Subheading")},
			match.List{},
			// Bad index
			match.NewIndex(4, match.NewAnyNode()),
		}
		_, err = query.QueryTree(tree, matcher)
		require.ErrorContains(t, err, "failed to match query: document[# Title][## Subheading].list did not have a [4].any")
	})
}
