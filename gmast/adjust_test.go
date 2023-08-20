package gmast_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/gmast"
)

func TestAppendChild(t *testing.T) {
	source := []byte("")

	newSegment, source := gmast.NewSegment("Hello, world!", source)

	list := gmast.AppendChild(ast.NewList('-'),
		gmast.AppendChild(ast.NewListItem(2),
			gmast.AppendChild(
				ast.NewTextBlock(),
				ast.NewTextSegment(newSegment),
			)))

	// Added the new text to the source
	require.Contains(t, string(source), "Hello, world!")
	fmt.Println(list)
	require.Equal(t, list.Kind(), ast.KindList)
}
