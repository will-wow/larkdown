package markdown_test

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown/gmast"
	"github.com/will-wow/larkdown/match"
	"github.com/will-wow/larkdown/query"
	"github.com/will-wow/larkdown/renderer/markdown"
)

var postMarkdown = `
# Markdown in Go

## Tags

#markdown #golang

In this essay I will explain...
`

func ExampleNewRenderer() {
	source := []byte(postMarkdown)
	// Preprocess the markdown into goldmark AST
	md := goldmark.New(
		// Parse hashtags to they can be matched against.
		goldmark.WithExtensions(
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
		),
	)
	doc := md.Parser().Parse(text.NewReader(source))

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

	// ====
	// Use larkdown renderer to render back to markdown
	// ====
	var newMarkdown bytes.Buffer
	err = markdown.NewNodeRenderer().Render(&newMarkdown, source, doc)
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
	// # Markdown in Go
	//
	// ## Tags
	//
	// #markdown #golang #testing
	//
	// In this essay I will explain...
}
