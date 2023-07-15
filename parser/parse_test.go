package parser_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/will-wow/larkdown/parser"
)

type List struct {
	Items []string `json:"items"`
}

type Ingredients struct {
	Buy    List `json:"buy"`
	OnHand List `json:"on-hand"`
}

type Recipe struct {
	Ingredients  Ingredients `json:"ingredients"`
	Instructions List        `json:"instructions"`
}

func TestParse(t *testing.T) {
	// t.Run("Parse", func(t *testing.T) {
	// 	mdSpec := spec.NewSpecDocument(spec.WithChildren([]spec.SpecNode{
	// 		spec.NewSpecHeading("my-heading", 1),
	// 	}))

	// 	source := []byte(md)

	// 	got, err := parser.MarkdownToTree(source, *mdSpec)
	// 	require.NoError(t, err)

	// 	childHeading, ok := got.SubHeadings["my-heading"]
	// 	require.True(t, ok)
	// 	require.Equal(t, 1, childHeading.Level)
	// 	require.Equal(t, "My Heading", string(childHeading.Text(source)))

	// 	subHeading1, ok := got.SubHeadings["my-heading"].SubHeadings["my-subheading"]
	// 	require.True(t, ok)
	// 	require.Equal(t, 2, subHeading1.Level)
	// 	require.Equal(t, "my-heading", subHeading1.Parent.Id)
	// 	require.Equal(t, "My Subheading", string(subHeading1.Text(source)))

	// 	subHeading2, ok := got.SubHeadings["my-heading"].SubHeadings["second-subheading"]
	// 	require.True(t, ok)
	// 	require.Equal(t, 2, subHeading2.Level)
	// 	require.Equal(t, "my-heading", subHeading2.Parent.Id)
	// 	require.Equal(t, "Second Subheading", string(subHeading2.Text(source)))
	// })

	t.Run("ParseRecipe", func(t *testing.T) {
		source := []byte(recipe)

		tree, err := parser.MarkdownToTree(source)
		require.NoError(t, err)

		query := []parser.NodeQuery{
			parser.BranchQuery{Level: 2, Name: []byte("Ingredients")},
			parser.BranchQuery{Level: 3, Name: []byte("Buy")},
			parser.IndexQuery{Index: 1, Query: parser.ListQuery{}},
		}

		match, err := parser.FindMatch(tree, query, source)
		require.NoError(t, err)

		list := &parser.ListParser{}
		parser.ParseResult(match, source, list)

		require.Equal(t, []string{"1 Medium Apple", "1 small-medium carrot", "1 banana", "2 eggs"}, list.Items)
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

var recipe = `
#baby-breakfast #baby #dairy-free #muffins

These Baby Led Weaning Muffins have no added sugar perfect for babies, toddlers, and kids. A Soft spongy style Healthy Baby Muffin with Apple Banana and Carrot.

## Ingredients

### Buy

- 1 Medium Apple
- 1 small-medium carrot
- 1 banana
- 2 eggs

### On hand

- 1 teaspoon vanilla essence
- 3 tbsp melted coconut oil (or butter)
- 1 1/4 cup flour
- 1.5 tsp baking powder

## Instructions

1. Peel, core, and dice apple
2. Peel and grate carrot, place apple and carrot in a pot with a little water, pop a lid on and simmer until apple soft. Usually 5-6 mins. Then drain apple/carrot mixture
3. While apple and carrot are cooking, place banana in a large bowl, mash with a fork
4. Add eggs, vanilla, and butter/oil
5. Puree the cooked apple and carrot, I use a stick blender
6. Add apple and carrot to the other wet ingredients
7. Beat these wet ingredients together with a hand-held beater, should become smooth, yellow and a little frothy
8. Add the flour and baking powder
9. Beat for a further 30-60 seconds to make a well-mixed batter
10. Portion into an oiled muffin tin (I use a non-stick silicon tray sprayed with oil) Mix makes 12 standard sized muffins or 30 mini muffins
11. Bake at 180 degrees Celsius for 15 mins (350 Fahrenheit) 15 min cook time is based on making mini muffins, the mix makes approx 30 mini muffins. If you are using a standard muffin tray and making approx 12 muffins the cook time will be longer, approx 25-30 mins
12. Cool
13. Serve
14. These muffins can be stored in an airtight container for 3 days, or they can be frozen.
`
