package larkdown

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown/gmast"
)

func DecodeListItems(node ast.Node, source []byte) (out []string, err error) {
	list, ok := node.(*ast.List)
	if !ok {
		return out, fmt.Errorf("expected list node")
	}

	gmast.ForEachListItem(list, source, func(item ast.Node, _ int) {
		out = append(out, string(item.Text(source)))
	})

	return out, nil
}

func DecodeText(node ast.Node, source []byte) (string, error) {
	return string(node.Text(source)), nil
}

func DecodeTag(node ast.Node, source []byte) (string, error) {
	tag, ok := node.(*hashtag.Node)
	if !ok {
		return "", fmt.Errorf("expected hashtag node, got %s", node.Text(source))
	}

	return string(tag.Tag), nil
}
