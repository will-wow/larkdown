package parser

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown/spec"
)

var attrNameID = []byte("id")

// Parse parses a markdown document into a Heading tree
func Parse(source []byte, spec spec.SpecDocument) (Heading, error) {
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

		child = child.NextSibling()
	}

	return tree, nil
}

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
