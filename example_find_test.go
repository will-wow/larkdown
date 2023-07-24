package larkdown_test

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/match"
)

var findMarkdown = `
# My Recipe

Here's a long story about making dinner.

## Tags

#dinner #chicken

## Ingredients

- Chicken
- Vegetables
- Salt
- Pepper
`

func ExampleFind() {
	source := []byte(findMarkdown)
	// Preprocess the markdown into goldmark AST
	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(source))

	// Set up a ingredientsQuery for the first list under ## Ingredients
	ingredientsQuery := []match.Node{
		match.Branch{Level: 1},
		match.Branch{Level: 2, Name: []byte("Ingredients")},
		match.Index{Index: 0, Node: match.List{}},
	}

	// Decode the list items into a slice of strings
	ingredients, err := larkdown.Find(doc, source, ingredientsQuery, larkdown.DecodeListItems)
	if err != nil {
		panic(fmt.Errorf("couldn't find an ingredients list: %w", err))
	}

	fmt.Println(ingredients)

	// Output:
	// [Chicken Vegetables Salt Pepper]
}
