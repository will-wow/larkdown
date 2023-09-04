# larkdown

Lock down your markdown.

Larkdown allows you to treat markdown files as a tree where headings are branches, to extract data from that tree, and then either update and render the tree back to markdown, or continue to render to HTML.

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

and then query that data structure to find a node. With a node you can then decode it into useful data like strings and slices of strings, or change it and re-save back to markdown.

Specially `larkdown` takes an AST generated from the excellent [goldmark](https://github.com/yuin/goldmark) library for parsing [Commonmark](https://commonmark.org) markdown, and lest you query, update, and re-render that AST. This makes it easy to take a markdown file, run it through Goldmark, query some structured data, and then either finish using Goldmark to render the file to HTML, or make some updates and save back to markdown.

## Motivation

This library acts as a test bed for an idea - markdown has the excellent property of being good for both machine and human reading. Therefor (with the right tooling) it should be possible to use Markdown files as an externally portable data store. You could author and edit data using one tool, and then serve it on the web (with editing capabilities) with another tool. And if you want to switch or give up a tool at some point, it's no problem - it's just markdown, you don't need to export or transform it.

Larkdown is an attempt to build that tooling.

## Usage

```bash
go get github.com/will-wow/larkdown
```

You can use larkdown to pull out important data about a file before sending it to a frontend to be rendered:

```go
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
	err = md.Renderer().Render(&recipe.Html, source, doc)
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
```

Or you can use it to update a markdown file in-place, and still render to HTML afterwards:

```go
package mdrender_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/gmast"
	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/mdfront"
	"github.com/will-wow/larkdown/mdrender"
	"github.com/will-wow/larkdown/query"
)

var postMarkdown = `# Markdown in Go

## Tags

#markdown #golang

In this essay I will explain...
`

type PostData struct {
	Slug string `yaml:"slug"`
}

// Matcher for a line of #tags under the heading ## Tags
var tagsQuery = []match.Node{
	match.Branch{Level: 2, Name: []byte("tags"), CaseInsensitive: true},
	match.NodeOfKind{Kind: ast.KindParagraph},
}

var titleQuery = []match.Node{match.Heading{Level: 1}}

func ExampleNewRenderer() {
	source := []byte(postMarkdown)
	// Preprocess the markdown into goldmark AST
	md := goldmark.New(
		goldmark.WithExtensions(
			// Parse hashtags to they can be matched against.
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			// Support frontmatter rendering.
			// This does nothing on its own, but sets up a place to render frontmatter to.
			&mdfront.Extender{},
		),
	)

	// Set up context for the metadata
	context := parser.NewContext()
	// Parse the markdown into an AST, with context
	doc := md.Parser().Parse(text.NewReader(source), parser.WithContext(context))

	// ====
	// Get the tags from the file
	// ====

	// Find the tags header to append to
	tagsLine, err := query.QueryOne(doc, source, tagsQuery)
	if err != nil {
		panic(fmt.Errorf("error finding tags heading: %w", err))
	}

	// ====
	// Edit the AST to add a new tag
	// ====

	// Create a new tag
	space, source := gmast.NewSpace(source)
	hashtag, source := gmast.NewHashtag("testing", source)

	// Append the new tag to the tags line
	gmast.AppendChild(tagsLine,
		space,
		hashtag,
	)

	// ====
	// Add a slug to the post's frontmatter.
	// ====

	// Find the title header to use as a slug
	// In practice, you might want to also pull this frontmatter out of an existing document
	// using goldmark-meta or goldmark-frontmatter.
	title, err := larkdown.Find(doc, source, titleQuery, larkdown.DecodeText)
	if err != nil {
		panic(fmt.Errorf("error finding title: %w", err))
	}

	// Slugify the title
	slug := strings.ReplaceAll(strings.ToLower(title), " ", "-")

	// Set up a struct to render the frontmatter
	data := &PostData{Slug: slug}

	// ====
	// Use larkdown renderer to render back to markdown
	// ====
	var newMarkdown bytes.Buffer
	// Here we set up the renderer outside the goldmark.New call, so you can use the normal
	// goldmark HTML renderer, and also render back to markdown.
	err = larkdown.NewNodeRenderer(
		// Pass the metaData to the renderer to render back to markdown
		mdrender.WithFrontmatter(data),
	).Render(&newMarkdown, source, doc)
	if err != nil {
		panic(fmt.Errorf("error rendering Markdown: %w", err))
	}

	// ====
	// Also render to HTML
	// ====
	var html bytes.Buffer
	err = md.Renderer().Render(&html, source, doc)
	if err != nil {
		panic(fmt.Errorf("error rendering HTML: %w", err))
	}

	// The new #testing tag is after the #golang tag in the HTML output
	fmt.Println("HTML:")
	fmt.Println(html.String())

	// The new #testing tag is after the #golang tag in the markdown
	fmt.Println("Markdown:")
	fmt.Println(newMarkdown.String())

	// Output:
	// HTML:
	// <h1>Markdown in Go</h1>
	// <h2>Tags</h2>
	// <p><span class="hashtag">#markdown</span> <span class="hashtag">#golang</span> <span class="hashtag">#testing</span></p>
	// <p>In this essay I will explain...</p>
	//
	// Markdown:
	// ---
	// slug: markdown-in-go
	// ---
	//
	// # Markdown in Go
	//
	// ## Tags
	//
	// #markdown #golang #testing
	//
	// In this essay I will explain...
}
```

## Roadmap

- [x] basic querying and unmarshaling of headings, lists, and text
- [x] make sure this works with extracting front matter
- [x] make sure this doesn't interfere with rendering the markdown to HTML with goldmark
- [x] tag matchers/decoders
- [x] handle finding multiple matches
- [x] generic matcher for any goldmark kind
- [x] basic markdown renderer
- [ ] Full markdown renderer
- [x] Basic markdown editing tools
- [ ] Move markdown editing tools
- [ ] options for recording extra debugging data for failed matches
- [ ] use options to support not setting a matcher or decoder
- [ ] handle decoding a table into a slice of structs
- [ ] handle a list of matchers for FindAll extractors
- [ ] matchers/decoders for more nodes:
  - [ ] codeblocks by language
  - [ ] tables with structured output
- [ ] add an "end on" option for branches, to end on the next subheading of a specific level
- [ ] nth instance matcher for queries like "the second list"
- [ ] query validator to make sure it even makes sense
- [ ] query syntax based on CSS selectors
- [ ] Update queries to fit with CSS selectors
- [ ] cli for selector queries
- [ ] generic unmarshaler into json
- [x] benchmark
- [x] more docs and tests

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

### Publish

```bash
VERSION=0.0.N task publish
```
