package gmast_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/gmast"
	"github.com/will-wow/larkdown/internal/test"
)

func TestForEachListItem(t *testing.T) {
	tree, source := test.TreeFromMd(t, `
	- foo
	- foo bar
	- baz
	`)

	list, ok := tree.FirstChild().(*ast.List)
	require.True(t, ok, "first child is not a list")

	items := []string{}
	gmast.ForEachListItem(list, source, func(item ast.Node, index int) {
		items = append(items, string(item.Text(source)))
	})

	require.Equal(t, []string{"foo", "foo bar", "baz"}, items)
}

func TestForEachChild(t *testing.T) {
	tree, source := test.TreeFromMd(t, `
	foo

	foo bar

	baz
	`)

	items := []string{}
	gmast.ForEachChild(tree, source, func(item ast.Node, index int) {
		items = append(items, string(item.Text(source)))
	})

	require.Equal(t, []string{"foo", "foo bar", "baz"}, items)
}

func TestWalkSiblingsUntil(t *testing.T) {
	tree, source := test.TreeFromMd(t, `
	foo

	bar

	stop
	`)

	items := []string{}
	err := gmast.WalkSiblingsUntil(tree.FirstChild(), func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if string(n.Text(source)) == "stop" {
			return ast.WalkStop, nil
		}

		items = append(items, n.Kind().String())
		return ast.WalkContinue, nil
	})
	require.NoError(t, err)

	require.Equal(t, []string{"Paragraph", "Text", "Paragraph", "Text"}, items)
}

func TestFindSibling(t *testing.T) {
	tree, source := test.TreeFromMd(t, `
	## Heading 2

	Body H2

	### Heading 3

	Body H3

	## Heading 2 Again

	Body H2.2
	`)

	secondH2 := gmast.FindSibling(tree.FirstChild(), func(node ast.Node) bool {
		heading, ok := node.(*ast.Heading)
		if !ok {
			return false
		}

		return heading.Level == 2
	})

	require.Equal(t, "Heading 2 Again", string(secondH2.Text(source)))
}

func TestLastChildOfHeading(t *testing.T) {
	tree, source := test.TreeFromMd(t, `
	## Heading 2

	Body H2

	### Heading 3

	Body H3

	## Heading 2 Again

	Body H2.2
	`)

	lastChild, err := gmast.LastChildOfHeading(tree.FirstChild())
	require.NoError(t, err, "error finding last child of heading 2")

	require.Equal(t, "Body H3", string(lastChild.Text(source)))
}
