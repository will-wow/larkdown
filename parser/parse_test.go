package parser_test

import (
	"testing"

	"github.com/will-wow/larkdown/parser"
)

func TestParse(t *testing.T) {
	t.Run("Parse", func(t *testing.T) {
		spec := parser.NewSpecDocument(parser.WithChildren([]parser.SpecNode{
			parser.NewSpecHeading("my-heading", 1),
		}))

		// want := parser.Document{parser.NewBaseNode(
		// 	[]parser.Node{
		// 		&parser.Heading{Id: "my-heading"},
		// 	},
		// )}

		source := []byte(md)

		got, err := parser.Parse(source, *spec)
		if err != nil {
			t.Fatalf("got error %v", err)
		}

		childHeading, ok := got.SubHeadings["my-heading"]
		if !ok {
			t.Fatal("my-heading not found")
		}
		if childHeading.Level != 1 {
			t.Fatalf("got level %d, want %d", childHeading.Level, 1)
		}
		if text := childHeading.Text(source); string(text) != "My Heading" {
			t.Fatalf("got text %s, want %s", string(text), "My Heading")
		}

		subHeading1, ok := got.SubHeadings["my-heading"].SubHeadings["my-subheading"]
		if !ok {
			t.Fatal("my-subheading not found")
		}
		if subHeading1.Level != 2 {
			t.Fatalf("got level %d, want %d", childHeading.Level, 2)
		}
		if subHeading1.Parent.Id != "my-heading" {
			t.Fatalf("got parent id %s, want %s", subHeading1.Parent.Id, "my-heading")
		}
		if text := subHeading1.Text(source); string(text) != "My Subheading" {
			t.Fatalf("got text %s, want %s", string(text), "My Subheading")
		}

		subHeading2, ok := got.SubHeadings["my-heading"].SubHeadings["second-subheading"]
		if !ok {
			t.Fatal("second-subheading not found")
		}
		if subHeading2.Level != 2 {
			t.Fatalf("got level %d, want %d", childHeading.Level, 2)
		}
		if subHeading2.Parent.Id != "my-heading" {
			t.Fatalf("got parent id %s, want %s", subHeading2.Parent.Id, "my-heading")
		}
		if text := subHeading2.Text(source); string(text) != "Second Subheading" {
			t.Fatalf("got text %s, want %s", string(text), "Second Subheading")
		}
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
