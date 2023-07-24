// Test helpers
package test

import (
	"os"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Construct a goldmark tree from a markdown string, removing leading whitespace.
func TreeFromMd(t testing.TB, markdown string, opts ...goldmark.Option) (tree ast.Node, source []byte) {
	t.Helper()

	source = []byte(dedent.Dedent(markdown))

	md := goldmark.New(opts...)
	return md.Parser().Parse(text.NewReader(source)), source
}

// Construct a goldmark tree from a file.
func TreeFromFile(t testing.TB, filePath string, opts ...goldmark.Option) (tree ast.Node, source []byte) {
	t.Helper()

	source, err := os.ReadFile(filePath)
	require.NoError(t, err)

	md := goldmark.New(opts...)
	return md.Parser().Parse(text.NewReader(source)), source
}
