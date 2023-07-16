package examples

import (
	"fmt"

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

func ParseFile(filename string) (results []string, err error) {
	// Preprocess the markdown into a tree where headings are branches.
	tree, err := larkdown.MarkdownToTree([]byte(md))
	if err != nil {
		return results, fmt.Errorf("couldn't parse markdown: %w", err)
	}

	// Set up a matcher for find your data in the tree.
	matcher := []match.Node{
		match.Branch{Level: 1},
		match.Branch{Level: 2, Name: []byte("Ingredients")},
		match.Index{Index: 1, Node: match.List{}},
	}

	// Set up a NodeUnmarshaler to parse and store the data you want
	list := &larkdown.StringList{}
	err = larkdown.Unmarshal(tree, matcher, list)
	if err != nil {
		return results, fmt.Errorf("couldn't find an ingredients list: %w", err)
	}

	// Returns []string{"Chicken", "Vegetables", "Salt", "Pepper"}
	return list.Items, nil
}
