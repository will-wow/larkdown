# larkdown

Lock down your markdown.

Larkdown allows you to treat markdown files as a tree where headings are branches, to extract data from that tree.

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

and then query that data structure to find a node, and decode that node into useful data like strings and slices of strings.

Specially `larkdown` takes an AST generated from the excellent [goldmark](https://github.com/yuin/goldmark) library for parsing [Commonmark](https://commonmark.org) markdown, and lest you query that AST. This makes it easy to take a markdown file, run it through Goldmark, query some structured data, and then finish using Goldmark to render the file to HTML.

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

```

## Roadmap

- [x] basic querying and unmarshaling of headings, lists, and text
- [x] make sure this works with extracting front matter
- [x] make sure this doesn't interfere with rendering the markdown to HTML with goldmark
- [x] tag matchers/decoders
- [x] handle finding multiple matches
- [x] generic matcher for any goldmark kind
- [ ] options for recording extra debugging data for failed matches
- [ ] handle decoding a table into a slice of structs
- [ ] system for wrapping a set of matchers into one, for things like "give me the first line of every list under this subheader"
- [ ] matchers/decoders for more nodes like codeblocks by language
- [ ] nth instance matcher for queries like "the second list"
- [ ] string-based query syntax, ie. `"['# heading']['## heading2'].list[1]"`
- [ ] generic unmarshaler into json
- [ ] cli for string-based queries and json return values
- [ ] more docs and tests

## Alternatives

- [markdown-to-json](https://github.com/njvack/markdown-to-json): Python-based library for parsing markdown into JSON with a similar nested style.

## Contributing

### Install task

This project uses [task](https://taskfile.dev) as its task runner.

```bash
# macos
brew install go-task/tap/go-task

# linux/wsl
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
```

Or follow the [installation instructions](https://taskfile.dev/installation/) for more options.

### List commands

For a list of all commands for this project, run

```bash
task --list
```

### Format

```bash
task fmt
```

### Lint

```bash
task lint
```

### Test

```bash
task test
```

### Make ready for a commit

```bash
# runs fmt lint test
task ready
```
