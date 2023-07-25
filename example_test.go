package larkdown_test

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/match"
)

var recipeMarkdown = `
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

type Recipe struct {
	Tags        []string
	Ingredients []string
	Html        bytes.Buffer
}

func Example() {
	source := []byte(recipeMarkdown)
	// Preprocess the markdown into goldmark AST
	md := goldmark.New(
		// Parse hashtags to they can be matched against.
		goldmark.WithExtensions(
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
		),
	)
	doc := md.Parser().Parse(text.NewReader(source))

	recipe := Recipe{}

	// ====
	// Get the ingredients from the list
	// ====

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
	recipe.Ingredients = ingredients

	// ====
	// Get the tags from the file
	// ====

	// Matcher for the tags header
	tagsQuery := []match.Node{
		match.Branch{Level: 2, Name: []byte("Tags")},
	}

	// Find all Tags under the tags header, and decode their contents into strings.
	tags, err := larkdown.FindAll(doc, source, tagsQuery, match.Tag{}, larkdown.DecodeTag)
	if err != nil {
		// This will not return an error if there are no tags, only if something else went wrong.
		panic(fmt.Errorf("error finding tags: %w", err))
	}
	recipe.Tags = tags

	// ====
	// Render the HTML
	// ====
	err = md.Convert(source, &recipe.Html)
	if err != nil {
		panic(fmt.Errorf("error rendering HTML: %w", err))
	}

	fmt.Println(recipe.Ingredients)
	fmt.Println(recipe.Tags)
	fmt.Println(recipe.Html.String())

	// Output:
	// [Chicken Vegetables Salt Pepper]
	// [dinner chicken]
	// <h1>My Recipe</h1>
	// <p>Here's a long story about making dinner.</p>
	// <h2>Tags</h2>
	// <p><span class="hashtag">#dinner</span> <span class="hashtag">#chicken</span></p>
	// <h2>Ingredients</h2>
	// <ul>
	// <li>Chicken</li>
	// <li>Vegetables</li>
	// <li>Salt</li>
	// <li>Pepper</li>
	// </ul>
}
