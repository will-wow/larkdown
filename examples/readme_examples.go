package examples

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/match"
)

var md = `
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

func ParseRecipe() (recipe Recipe, err error) {
	source := []byte(md)
	// Preprocess the markdown into goldmark AST
	md := goldmark.New(
		// Parse hashtags to they can be matched against.
		goldmark.WithExtensions(
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
		),
	)
	doc := md.Parser().Parse(text.NewReader(source))

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
		return recipe, fmt.Errorf("couldn't find an ingredients list: %w", err)
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
		return recipe, fmt.Errorf("error finding tags: %w", err)
	}
	recipe.Tags = tags

	// ====
	// Render the HTML
	// ====
	md.Convert(source, &recipe.Html)

	// Return all recipe data
	return recipe, nil
}
