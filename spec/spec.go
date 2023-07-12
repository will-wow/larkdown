package spec

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
