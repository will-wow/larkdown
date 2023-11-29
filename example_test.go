package larkdown_test

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
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

## Comments

| Name      | Comment    |
| --------- | ---------- |
| Alice     | It's good! |
| Bob       | It's bad   |

`

type Comment struct {
	Name    string
	Comment string
}

func (c Comment) String() string {
	return fmt.Sprintf("%s: %s", c.Name, c.Comment)
}

type Recipe struct {
	Tags        []string
	Ingredients []string
	Comments    []Comment
	Html        bytes.Buffer
}

func Example() {
	source := []byte(recipeMarkdown)
	// Preprocess the markdown into goldmark AST
	md := goldmark.New(
		// Parse hashtags to they can be matched against.
		goldmark.WithExtensions(
			extension.Table,
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

	tableQuery := []match.Node{
		match.Branch{Level: 2, Name: []byte("Comments")},
		match.Table{},
	}

	// ====
	// Get the comments from a table
	// ====

	// Get data from the comments table
	commentsTable, err := larkdown.Find(doc, source, tableQuery, larkdown.DecodeTableToMap)
	if err != nil {
		panic(fmt.Errorf("error finding comments: %w", err))
	}
	for _, comment := range commentsTable {
		recipe.Comments = append(recipe.Comments, Comment{
			Name:    comment["Name"],
			Comment: comment["Comment"],
		})
	}

	fmt.Println(recipe.Ingredients)
	fmt.Println(recipe.Tags)
	fmt.Println(recipe.Comments)

	// Output:
	// [Chicken Vegetables Salt Pepper]
	// [dinner chicken]
	// [Alice: It's good! Bob: It's bad]
}
