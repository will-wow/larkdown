package mdrender_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown"
	"github.com/will-wow/larkdown/gmast"
	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/mdrender"
	"github.com/will-wow/larkdown/query"
)

var postMarkdown = `---
slug: markdown-in-go
---

# Markdown in Go

## Tags

#markdown #golang

In this essay I will explain...
`

type PostData struct {
	Slug string `yaml:"slug"`
}

func ExampleNewRenderer() {
	source := []byte(postMarkdown)
	// Preprocess the markdown into goldmark AST
	md := goldmark.New(
		goldmark.WithExtensions(
			// Parse hashtags to they can be matched against.
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			meta.Meta,
			// Support frontmatter rendering.
			&mdrender.Extender{},
		),
	)

	// Set up context for the metadata
	context := parser.NewContext()
	// Parse the markdown into an AST, with context
	doc := md.Parser().Parse(text.NewReader(source), parser.WithContext(context))

	// ====
	// Get the tags from the file
	// ====

	// Matcher for the tags line
	tagsQuery := []match.Node{
		match.Branch{Level: 2, Name: []byte("tags"), CaseInsensitive: true},
		match.NodeOfKind{Kind: ast.KindParagraph},
	}

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

	// Set up a struct to pull in the frontmatter
	data := &PostData{}

	var slug string
	// Populate the meta struct with the slug
	contextSlug := meta.Get(context)["slug"]
	if contextSlug == nil {
		// Find the title header to use as a slug
		title, err := larkdown.Find(doc, source, []match.Node{match.Heading{Level: 1}}, larkdown.DecodeText)
		if err != nil {
			panic(fmt.Errorf("error finding title: %w", err))
		}

		// Slugify the title
		slug = strings.ReplaceAll(strings.ToLower(title), " ", "-")
	} else {
		// Otherwise keep the slug from the frontmatter
		slug = fmt.Sprint(contextSlug)
	}
	// Store the slug for re-rendering
	data.Slug = slug

	// ====
	// Use larkdown renderer to render back to markdown
	// ====
	var newMarkdown bytes.Buffer
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
