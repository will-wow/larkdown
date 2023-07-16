package larkdown

import (
	"fmt"

	"github.com/yuin/goldmark/ast"

	"github.com/will-wow/larkdown/gmast"
)

// Implement this interface to unmarshal a goldmark AST node.
type NodeUnmarshaler interface {
	UnmarshalNode(node ast.Node, source []byte) error
}

// Unmarshals a goldmark List node's items into a slice of strings.
type StringList struct {
	Items []string
}

func (u *StringList) UnmarshalNode(node ast.Node, source []byte) error {
	list, ok := node.(*ast.List)
	if !ok {
		return fmt.Errorf("expected list node")
	}

	gmast.ForEachListItem(list, source, func(item ast.Node, _ int) {
		u.Items = append(u.Items, string(item.Text(source)))
	})

	return nil
}

// Unmarshals any goldmark node's text into a string.
type NodeText struct {
	Text string
}

func (u *NodeText) UnmarshalNode(node ast.Node, source []byte) error {
	u.Text = string(node.Text(source))
	return nil
}
