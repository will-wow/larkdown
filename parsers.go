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

func (p *StringList) UnmarshalNode(node ast.Node, source []byte) error {
	list, ok := node.(*ast.List)
	if !ok {
		return fmt.Errorf("expected list node")
	}

	gmast.ForEachListItem(list, source, func(item ast.Node, _ int) {
		p.Items = append(p.Items, string(item.Text(source)))
	})

	return nil
}
