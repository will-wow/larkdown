# larkdown

Lock down your markdown.

This package allows you to treat markdown files as a tree where headings are branches, and then extract data from that trees.

It lets you treat this:

```markdown
# Title

## Subheading

### Sub-subheading

- a list
- of things

## Another subheading

Some content
```

like this

```json
{
  "Title": [
    "# Title",
    {
      "Subheading": [
        "## Subheading",
        {
          "Sub-subheading": [
            "### Sub-subheading",
            {
              "list": ["a list", "of things"]
            }
          ]
        }
      ]
    },
    {
      "Another subheading": ["## Another subheading", "some content"]
    }
  ]
}
```

and then query that data structure to find a node, and unmarshal that node into useful data like strings nad slices of strings.

`larkdown` makes heavy use of the excellent [goldmark](https://github.com/yuin/goldmark) library for parsing [Commonmark](https://commonmark.org) markdown into a structure that is easy to work with.

## Motivation

I do a lot of cooking, and I keep my recipes as markdown files edited through [Obsidian](https://obsidian.md) for ease of authoring and portability. I wanted to write some tooling to make it easier to build grocery lists for the week, and wanted to take advantage of the fact that my recipes were already regularly structured, with an `## Ingredients` heading followed by a list. This library lets me pull out that ingredient data, so I can make a shopping list and get on with my weekend.

## Usage

```bash
go get github.com/will-wow/larkdown
```

```go
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
```

## Roadmap

- [x] basic querying and unmarshaling of headings, lists, and text
- [ ] matchers/unmarshalers for more nodes like codeblocks by language
- [ ] nth instance matcher for queries like "the second list"
- [ ] string-based query syntax, ie. `"['# heading']['## heading2'].list[1]"`
- [ ] generic unmarshaler into json
- [ ] cli for string-based queries and json return values
- [ ] more docs and tests

## Alternatives

- [markdown-to-json](https://github.com/njvack/markdown-to-json): Python-based library for parsing markdown into JSON with a similar nested style.

## Contributing

### Format

```bash
make fmt
```

### Lint

```bash
make lint
```

### Test

```bash
make test
```

### Make ready for a commit

```bash
# runs fmt lint test
make ready
```
