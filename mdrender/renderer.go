package mdrender

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/hashtag"
	"gopkg.in/yaml.v3"

	"github.com/will-wow/larkdown/mdfront"
)

// A Config struct has configuration for the markdown renderer.
type Config struct {
	Writer      Writer // Writer is a writer used to write textual contents.
	Frontmatter any    // Frontmatter is a frontmatter struct to be rendered to yaml.
}

// NewConfig returns a new Config with defaults.
func NewConfig() Config {
	return Config{
		Writer: DefaultWriter,
	}
}

// SetOption implements renderer.NodeRenderer.SetOption.
func (c *Config) SetOption(name renderer.OptionName, value interface{}) {
	if name == optTextWriter {
		c.Writer, _ = value.(Writer)
	}
}

// An Option interface sets options for HTML based renderers.
type Option interface {
	SetMarkdownOption(*Config)
}

// TextWriter is an option name used in WithWriter.
const optTextWriter renderer.OptionName = "Writer"

type withWriter struct {
	value Writer
}

func (o *withWriter) SetConfig(c *renderer.Config) {
	c.Options[optTextWriter] = o.value
}

func (o *withWriter) SetMarkdownOption(c *Config) {
	c.Writer = o.value
}

// WithWriter is a functional option that allow you to set the given writer to
// the renderer.
func WithWriter(writer Writer) interface {
	renderer.Option
	Option
} {
	return &withWriter{writer}
}

type withFrontmatter struct {
	data any
}

var _ Option = (*withFrontmatter)(nil)

func (o *withFrontmatter) SetMarkdownOption(c *Config) {
	c.Frontmatter = o.data
}

// WithFrontmatter records a frontmatter struct for repopulating the frontmatter.
func WithFrontmatter(data any) Option {
	return &withFrontmatter{data: data}
}

// A Renderer struct is an implementation of renderer.NodeRenderer that renders
// nodes as Markdown.
type Renderer struct {
	Config
}

var _ renderer.NodeRenderer = &Renderer{}

// NewRenderer returns a new Renderer with given options.
func NewRenderer(opts ...Option) renderer.NodeRenderer {
	r := &Renderer{
		Config: NewConfig(),
	}

	for _, opt := range opts {
		opt.SetMarkdownOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs .
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks

	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)

	// inlines

	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)
	reg.Register(hashtag.Kind, r.renderHashtag)

	// Frontmatter
	reg.Register(mdfront.Kind, r.renderFrontmatter)
}

func (r *Renderer) writeLines(w util.BufWriter, source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		r.Writer.Write(w, line.Value(source))
	}
}

func (r *Renderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// nothing to do
	return ast.WalkContinue, nil
}

func (r *Renderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, _ := node.(*ast.Heading)
	if entering {
		for i := 0; i < n.Level; i++ {
			_ = w.WriteByte('#')
		}
		_ = w.WriteByte(' ')
	} else {
		_, _ = w.WriteString("\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderBlockquote(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// For each paragraph:
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			// For each line in the paragraph:
			l := c.Lines().Len()
			for i := 0; i < l; i++ {
				// Write the line with a > prefix
				line := c.Lines().At(i)
				_, _ = w.WriteString("> ")
				_, _ = w.Write(line.Value(source))
			}

			if c.NextSibling() != nil {
				// Extra blank > line between paragraphs of a quote
				_, _ = w.WriteString("\n>\n")
			} else {
				_ = w.WriteByte('\n')
			}
		}

		return ast.WalkSkipChildren, nil
	} else {
		_ = w.WriteByte('\n')
		return ast.WalkContinue, nil
	}
}

func (r *Renderer) renderCodeBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("```\n")
		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString("```\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, _ := node.(*ast.FencedCodeBlock)
	if entering {
		_, _ = w.WriteString("```")
		language := n.Language(source)
		if language != nil {
			r.Writer.Write(w, language)
		}
		_, _ = w.WriteString("\n")
		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString("```\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, _ := node.(*ast.HTMLBlock)
	if entering {
		l := n.Lines().Len()
		for i := 0; i < l; i++ {
			line := n.Lines().At(i)
			r.Writer.Write(w, line.Value(source))
		}
	} else {
		if n.HasClosure() {
			closure := n.ClosureLine
			r.Writer.Write(w, closure.Value(source))
		}
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering && node.Parent().Kind() == ast.KindDocument && node.NextSibling() != nil {
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func indentListItemChild(w util.BufWriter, node ast.Node) {
	indent := 0
	for node.Parent() != nil {
		list, ok := node.Parent().(*ast.List)
		if ok {
			if list.IsOrdered() {
				indent += 3
			} else {
				indent += 2
			}
		}
		node = node.Parent()
	}

	for i := 0; i < indent; i++ {
		_ = w.WriteByte(' ')
	}
}

func (r *Renderer) renderListItem(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	list, _ := n.Parent().(*ast.List)

	if entering {
		ordered := list.IsOrdered()

		// Add indent for each parent list.
		indentListItemChild(w, list)

		if ordered {
			_ = w.WriteByte('1')
			_ = w.WriteByte(list.Marker)
		} else {
			_ = w.WriteByte(list.Marker)
		}

		// Add a space after the marker, before the children.
		if n.HasChildren() {
			_ = w.WriteByte(' ')
		}
	} else if listItemLineBreak(n, list.IsTight) {
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func listItemLineBreak(n ast.Node, isTight bool) bool {
	// Skip adding a newline if the last child is a list (since it already does that).
	lc := n.LastChild()
	if lc != nil && lc.Kind() == ast.KindList {
		return false
	}

	if !isTight {
		nc := n.NextSibling()
		if nc == nil || !nc.HasBlankPreviousLines() {
			return false
		}
	}

	return true
}

func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// If the paragraph is a child of a list item, indent it to match the list item.
		if n.Parent().Kind() == ast.KindListItem && n.Parent().FirstChild() != n {
			indentListItemChild(w, n)
		}
	} else {
		if n.Parent().LastChild() == n {
			_ = w.WriteByte('\n')
		} else {
			_, _ = w.WriteString("\n\n")
		}
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) renderTextBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if n.NextSibling() != nil && n.FirstChild() != nil {
			_ = w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderThematicBreak(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString("---\n\n")
	return ast.WalkContinue, nil
}

func (r *Renderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, _ := node.(*ast.AutoLink)
	if !entering {
		return ast.WalkContinue, nil
	}

	url := n.URL(source)

	_, _ = w.WriteString("<")
	// TODO: Labels?
	// w.Write(label)
	_, _ = w.Write(url)
	_, _ = w.WriteString(">")

	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	// TODO: somehow get the number of backticks, which isn't in the AST.
	// Write a ` on entering and leaving
	_ = w.WriteByte('`')

	return ast.WalkContinue, nil
}

func (r *Renderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, _ := node.(*ast.Emphasis)
	tag := "_"
	// TODO: Handle nested emphasis
	if n.Level == 2 {
		tag = "**"
	}

	// Write * on entering and leaving
	_, _ = w.WriteString(tag)
	return ast.WalkContinue, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, _ := node.(*ast.Link)

	if entering {
		_, _ = w.WriteString("[")
	} else {
		_, _ = w.WriteString("]")
		_, _ = w.WriteString("(")
		_, _ = w.Write(n.Destination)
		_, _ = w.WriteString(")")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n, _ := node.(*ast.Image)

	_, _ = w.WriteString("![")
	_, _ = w.Write(n.Text(source))
	_, _ = w.WriteString("](")
	_, _ = w.Write(n.Destination)

	if n.Title != nil {
		// If there's a title, write it in quotes.
		_, _ = w.WriteString(` "`)
		_, _ = w.Write(n.Title)
		_ = w.WriteByte('"')
	}
	_ = w.WriteByte(')')
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}

	n, _ := node.(*ast.RawHTML)
	l := n.Segments.Len()
	for i := 0; i < l; i++ {
		segment := n.Segments.At(i)
		_, _ = w.Write(segment.Value(source))
	}

	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n, _ := node.(*ast.Text)
	segment := n.Segment
	if n.IsRaw() {
		r.Writer.Write(w, segment.Value(source))
	} else {
		value := segment.Value(source)
		r.Writer.Write(w, value)
		if n.HardLineBreak() || n.SoftLineBreak() {
			_ = w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n, _ := node.(*ast.String)
	_, _ = w.Write(n.Value)
	return ast.WalkContinue, nil
}

func (r *Renderer) renderHashtag(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// No-op, since the child textSegment will have the #tag contents
	return ast.WalkContinue, nil
}

// renderFrontmatter renders the configured frontmatter back to yaml.
func (r *Renderer) renderFrontmatter(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	if r.Config.Frontmatter == nil {
		return ast.WalkContinue, nil
	}

	data, err := yaml.Marshal(r.Config.Frontmatter)
	if err != nil {
		return ast.WalkContinue, nil
	}

	// Print the frontmatter
	_, _ = w.WriteString("---\n")
	_, _ = w.Write(data)
	_, _ = w.WriteString("---\n\n")

	return ast.WalkContinue, nil
}

// A Writer interface writes textual contents to a writer.
type Writer interface {
	// Write writes the given source to writer with resolving references and unescaping
	// backslash escaped characters.
	Write(writer util.BufWriter, source []byte)
}

type defaultWriter struct{}

// NewWriter returns a new Writer.
func NewWriter() Writer {
	w := &defaultWriter{}
	return w
}

func (d *defaultWriter) Write(writer util.BufWriter, source []byte) {
	_, _ = writer.Write(source)
}

// DefaultWriter is a default instance of the Writer.
var DefaultWriter = NewWriter()
