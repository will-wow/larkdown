package preprocess

import "github.com/yuin/goldmark/ast"

// Represents a branch in the document tree, triggered by a heading.
// This is also implements goldmark/ast.Node, so it can be used with AST tools.
type TreeBranch struct {
	ast.BaseInline
	TreeParent *TreeBranch
	Level      int
}

func (n *TreeBranch) Dump(source []byte, level int) {
}

var KindTreeBranch = ast.NewNodeKind("TreeBranch")

func (n *TreeBranch) Kind() ast.NodeKind {
	return KindTreeBranch
}

// A tree branch for the root of the document.
func newTreeBranchRoot() *TreeBranch {
	return &TreeBranch{
		TreeParent: nil,
		Level:      0,
		BaseInline: ast.BaseInline{},
	}
}

func newTreeBranch(heading *ast.Heading, parent *TreeBranch) *TreeBranch {
	if heading == nil {
		panic("heading cannot be nil")
	}

	headingContents := *heading
	threeBranch := &TreeBranch{
		TreeParent: parent,
		Level:      headingContents.Level,
		BaseInline: ast.BaseInline{},
	}

	threeBranch.AppendChild(threeBranch, heading)

	return threeBranch
}
