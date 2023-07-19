package extension

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/will-wow/larkdown/preprocess"
)

type LarkdownASTTransformer struct{}

func NewLarkdownTransformer() parser.ASTTransformer {
	return &LarkdownASTTransformer{}
}

func (a *LarkdownASTTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	preprocess.GoldmarkToTree(node, reader.Source())
}

type LarkdownHTMLRenderer struct{}

func NewLarkdownHTMLRenderer() renderer.NodeRenderer {
	return &LarkdownHTMLRenderer{}
}

func (r *LarkdownHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(preprocess.KindTreeBranch, r.renderTreeBranch)
}

func (r *LarkdownHTMLRenderer) renderTreeBranch(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	// TODO: Implement?
	return ast.WalkContinue, nil
}

type LarkdownExtension struct{}

func NewLarkdownExtension() goldmark.Extender {
	return &LarkdownExtension{}
}

func (e *LarkdownExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithASTTransformers(util.Prioritized(
		NewLarkdownTransformer(), 200,
	)))

	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewLarkdownHTMLRenderer(), 500),
	))

}
