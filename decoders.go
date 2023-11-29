package larkdown

import (
	"fmt"
	"strconv"

	"github.com/yuin/goldmark/ast"
	extension_ast "github.com/yuin/goldmark/extension/ast"
	"go.abhg.dev/goldmark/hashtag"

	"github.com/will-wow/larkdown/gmast"
)

// Decode an ast.List into a slice of strings for each item
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

// Decode all the text inside any node
func DecodeText(node ast.Node, source []byte) (string, error) {
	return string(node.Text(source)), nil
}

// Decode a #tag parsed by go.abhg.dev/goldmark/hashtag into a string.
// Only the text content of the tag is returned, not the # prefix.
func DecodeTag(node ast.Node, source []byte) (string, error) {
	tag, ok := node.(*hashtag.Node)
	if !ok {
		return "", fmt.Errorf("expected hashtag node, got %s", node.Text(source))
	}

	return string(tag.Tag), nil
}

// DecodeTableToMap decodes a table node into a slice of maps of column names to column string values.
func DecodeTableToMap(node ast.Node, source []byte) ([]map[string]string, error) {
	rows := []map[string]string{}

	table, ok := node.(*extension_ast.Table)
	if !ok {
		return rows, fmt.Errorf("expected table node, got %s", node.Text(source))
	}

	headerRow := table.FirstChild()
	if headerRow == nil {
		return rows, nil
	}

	// Get rows
	headers := []string{}
	gmast.ForEachChild(table, source, func(row ast.Node, index int) {
		// Record the headers from the first row.
		if index == 0 && row.Kind() == extension_ast.KindTableHeader {
			gmast.ForEachChild(row, source, func(cell ast.Node, index int) {
				headers = append(headers, string(cell.Text(source)))
			})
			return
		}

		// Decode a row
		decodedRow := map[string]string{}
		gmast.ForEachChild(row, source, func(cell ast.Node, index int) {
			header := headers[index]
			if header == "" {
				// If there's no header, then just use the column index
				header = strconv.Itoa(index)
			}

			// Associate each cell with its header
			decodedRow[header] = string(cell.Text(source))
		})
		rows = append(rows, decodedRow)
	})

	return rows, nil
}
