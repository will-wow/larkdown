package examples

import (
	"fmt"
	"os"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/match"
)

func ParseFile(filename string) (results []string, err error) {
	source, err := os.ReadFile(filename)
	if err != nil {
		return results, fmt.Errorf("couldn't open file: %w", err)
	}

	// Preprocess the markdown into a tree where headings are branches.
	tree, err := larkdown.MarkdownToTree(source)
	if err != nil {
		return results, fmt.Errorf("couldn't parse markdown: %w", err)
	}

	// Set up a matcher for find your data in the tree.
	matcher := []match.Node{
		match.Branch{Level: 2, Name: []byte("Ingredients")},
		match.Branch{Level: 3, Name: []byte("Buy")},
		match.Index{Index: 1, Node: match.List{}},
	}

	// Set up a NodeUnmarshaler to parse and store the data you want
	list := &larkdown.StringList{}
	err = larkdown.Unmarshal(tree, matcher, list)
	if err != nil {
		return results, fmt.Errorf("couldn't find an ingredients list: %w", err)
	}

	// Returns []string{"1 Medium Apple", "1 small-medium carrot", "1 banana", "2 eggs"}
	return list.Items, nil
}
