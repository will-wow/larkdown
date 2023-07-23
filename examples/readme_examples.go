package examples

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/match"
)

var md = `
# My Recipe

Here's a long story about making dinner.

## Ingredients

- Chicken
- Vegetables
- Salt
- Pepper
`

func ParseFile() (results []string, err error) {
	source := []byte(md)
	// Preprocess the markdown into a tree where headings are branches.
	md := goldmark.New(
	// goldmark.WithExtensions(extension.NewLarkdownExtension()),
	)
	doc := md.Parser().Parse(text.NewReader(source))

	// Set up a matcher for find your data in the tree.
	matcher := []match.Node{
		match.Branch{Level: 1},
		match.Branch{Level: 2, Name: []byte("Ingredients")},
		match.Index{Index: 0, Node: match.List{}},
	}

	// Set up a NodeUnmarshaler to parse and store the data you want
	list, err := larkdown.Find(doc, source, matcher, larkdown.DecodeListItems)
	if err != nil {
		return results, fmt.Errorf("couldn't find an ingredients list: %w", err)
	}

	// Returns []string{"Chicken", "Vegetables", "Salt", "Pepper"}
	return list, nil
}
