package parser

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"
)

var attrNameID = []byte("id")

type SpecNode interface {
	setChildren(nodes []SpecNode)
	setOptional(optional bool)
}

type BaseSpecNode struct {
	optional bool
	// remove from final output
	// remove bool

	children []SpecNode
}

func (n *BaseSpecNode) setChildren(nodes []SpecNode) {
	n.children = nodes
}

func (n *BaseSpecNode) setOptional(optional bool) {
	n.optional = optional
}

type SpecOption func(n SpecNode)

func WithChildren(nodes []SpecNode) SpecOption {
	return func(n SpecNode) {
		n.setChildren(nodes)
	}
}

func WithOptional() SpecOption {
	return func(n SpecNode) {
		n.setOptional(true)
	}
}

type SpecDocument struct {
	BaseSpecNode
}

func applyOpts(n SpecNode, opts []SpecOption) {
	for _, opt := range opts {
		opt(n)

	}
}

func NewSpecDocument(opts ...SpecOption) *SpecDocument {
	x := SpecDocument{}

	applyOpts(&x, opts)

	return &x
}

type SpecHeading struct {
	BaseSpecNode
	level int
	id    string
}

func NewSpecHeading(id string, level int, opts ...SpecOption) *SpecHeading {
	x := SpecHeading{}

	applyOpts(&x, opts)

	return &x

}

type SpecList struct {
	BaseSpecNode
}

func NewSpecList(opts ...SpecOption) *SpecList {
	x := SpecList{}
	applyOpts(&x, opts)
	return &x
}

func Parse(source []byte, spec SpecDocument) (Heading, error) {
	md := goldmark.New(
		goldmark.WithExtensions(&hashtag.Extender{}),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	doc := md.Parser().Parse(text.NewReader(source))

	child := doc.FirstChild()
	tree := Heading{
		SubHeadings: make(map[string]*Heading),
	}
	var activeTreeNode *Heading
	activeTreeNode = &tree
	for {
		// If we are at the end of the document, break out of the loop
		if child == nil {
			break
		}

		switch node := child.(type) {
		case *ast.Heading:
			id, err := getHeadingId(*node, source)
			if err != nil {
				return Heading{}, err
			}

			// If the heading is at the same level as the active heading,
			// The parent is done and this is a sibling.
			if node.Level <= activeTreeNode.Level {
				activeTreeNode = activeTreeNode.Parent
			}

			newHeading := &Heading{
				Id:     id,
				Parent: activeTreeNode,
				Name:   string(node.Text(source)),
				Level:  node.Level,
				BaseNode: BaseNode{
					astNode: node,
				},
				SubHeadings: make(map[string]*Heading),
			}

			// Index the heading in the parent's subheadings
			activeTreeNode.SubHeadings[id] = newHeading
			// Record the heading in its parent's children
			activeTreeNode.children = append(activeTreeNode.children, newHeading)
			// Make this heading the active heading
			activeTreeNode = newHeading
		case *ast.Paragraph:
			newNode := &Paragraph{BaseNode: BaseNode{astNode: node}}
			// Record the heading in its parent's children
			activeTreeNode.children = append(activeTreeNode.children, newNode)
		default:
			return Heading{}, fmt.Errorf("unknown node type %T", node)
		}

		// if isNode(child) {
		// 	return child
		// }

		// if found := findNode(child, isNode); found != nil {
		// 	return found
		// }

		child = child.NextSibling()
	}

	return tree, nil
}

// func main() {
// 	spec := NewSpecDocument(WithChildren([]SpecNode{
// 		NewSpecHeading("ingredients", 2, WithChildren(
// 			[]SpecNode{
// 				NewSpecHeading("buy", 3, WithChildren(
// 					[]SpecNode{
// 						NewSpecList(),
// 					},
// 				)),
// 				NewSpecHeading("on-hand", 3, WithChildren(
// 					[]SpecNode{
// 						NewSpecList(),
// 					},
// 				)),
// 			})),
// 		NewSpecHeading("instructions", 2, WithChildren(
// 			[]SpecNode{
// 				NewSpecList(),
// 			},
// 		)),
// 		NewSpecHeading("notes", 2, WithOptional()),
// 	}))

// 	md := goldmark.New(
// 		goldmark.WithExtensions(&hashtag.Extender{}),
// 		goldmark.WithParserOptions(
// 			parser.WithAutoHeadingID(),
// 		),
// 	)

// 	source := `
// # hello hi hello

// #foo

// ## world

// foo
// bar

// baz

// boo

// - alpha
// - beta

// 1. one
// 2. two

// world!
// 	`

// 	doc := md.Parser().Parse(text.NewReader([]byte(source)))

// 	// activeSpecNode := spec

// 	// ast := Document{}

// 	child := doc.FirstChild()
// 	for {
// 		if child == nil {
// 			break
// 		}

// 		// if isNode(child) {
// 		// 	return child
// 		// }

// 		// if found := findNode(child, isNode); found != nil {
// 		// 	return found
// 		// }

// 		child = child.NextSibling()
// 	}

// 	// var child in go

// 	// variable in go

// 	// var child ast.Node
// 	// for i := 0; i < doc.ChildCount(); i++ {
// 	// 	if i == 0 {
// 	// 		child = doc.FirstChild()
// 	// 	} else {
// 	// 		child = child.NextSibling()
// 	// 	}

// 	// 	if v, ok := (child.(*ast.Heading)); ok {
// 	// 		fmt.Println("h", string(v.Text(src)))
// 	// 		idAttr, valid := v.Attribute([]byte(optAutoHeadingID))
// 	// 		if !valid {
// 	// 			continue
// 	// 		}
// 	// 		idBytes, ok := idAttr.([]byte)
// 	// 		if !ok {
// 	// 			continue
// 	// 		}

// 	// 	} else {

// 	// 		fmt.Println("p", string(child.Text(src)))
// 	// 	}

// 	// }

// }

// Find the first node that matches the given predicate.
// Using a depth-first search.
// func findNode(node ast.Node, isNode func(n ast.Node) bool) ast.Node {
// 	child := node.FirstChild()
// 	for {
// 		if child == nil {
// 			return nil
// 		}
// 		if isNode(child) {
// 			return child
// 		}

// 		if found := findNode(child, isNode); found != nil {
// 			return found
// 		}

// 		child = child.NextSibling()
// 	}

// }

// type Document struct {
// 	Children []Node
// }

type Node interface {
	Children() []Node
	// FindHeading(id string, level int) (Node, bool)
	// GetNthChild(n int) (Node, bool)
}

type BaseNode struct {
	children []Node
	astNode  ast.Node
}

func (n *BaseNode) Children() []Node {
	return n.children
}

func (n *BaseNode) Lines() *text.Segments {
	return n.astNode.Lines()
}

func (n *BaseNode) Text(source []byte) []byte {
	return n.astNode.Text(source)
}

type Document struct {
	BaseNode
}

type Heading struct {
	BaseNode
	Id          string
	Name        string
	Level       int
	Parent      *Heading
	SubHeadings map[string]*Heading
}

type Paragraph struct {
	BaseNode
}

// func (h *Heading) FindHeading(id string, level int) (Node, bool) {
// 	return nil, false
// }

// type HeadingGroup struct {
// 	name     string
// 	level    int
// 	children map[string]Group
// }

func getHeadingId(node ast.Heading, source []byte) (id string, err error) {
	idAttr, valid := node.Attribute(attrNameID)
	if !valid {
		return id, fmt.Errorf("heading %s has no id", node.Text(source))
	}
	idBytes, ok := idAttr.([]byte)
	if !ok {
		return id, fmt.Errorf("heading %s id cannot be decoded", node.Text(source))
	}
	return string(idBytes), nil
}
