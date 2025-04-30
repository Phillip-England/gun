package parser

type Element interface {
	GetTagName() string
	GetAttiribute(name string) (Attribute, bool)
	GetAttributes() []Attribute
}