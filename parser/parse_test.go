package parser_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/will-wow/larkdown/spec"

	"github.com/will-wow/larkdown/parser"
)

func TestParse(t *testing.T) {
	t.Run("Parse", func(t *testing.T) {
		spec := spec.NewSpecDocument(spec.WithChildren([]spec.SpecNode{
			spec.NewSpecHeading("my-heading", 1),
		}))

		source := []byte(md)

		got, err := parser.Parse(source, *spec)
		require.NoError(t, err)

		childHeading, ok := got.SubHeadings["my-heading"]
		require.True(t, ok)
		require.Equal(t, 1, childHeading.Level)
		require.Equal(t, "My Heading", string(childHeading.Text(source)))

		subHeading1, ok := got.SubHeadings["my-heading"].SubHeadings["my-subheading"]
		require.True(t, ok)
		require.Equal(t, 2, subHeading1.Level)
		require.Equal(t, "my-heading", subHeading1.Parent.Id)
		require.Equal(t, "My Subheading", string(subHeading1.Text(source)))

		subHeading2, ok := got.SubHeadings["my-heading"].SubHeadings["second-subheading"]
		require.True(t, ok)
		require.Equal(t, 2, subHeading2.Level)
		require.Equal(t, "my-heading", subHeading2.Parent.Id)
		require.Equal(t, "Second Subheading", string(subHeading2.Text(source)))
	})
}

var md = `
# My Heading

Big stuff

## My Subheading

Medium stuff

## Second Subheading

More medium stuff
`
