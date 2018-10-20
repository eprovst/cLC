package lamcalc

import (
	"strconv"
	"strings"
)

// Serialize returns the application as a De Bruijn index representation
func (lx Appl) Serialize() string {
	builder := strings.Builder{}

	lx.serialize(&builder)
	return builder.String()
}

// Serialize returns the lambda abstraction as a De Bruijn index representation
func (la Abst) Serialize() string {
	builder := strings.Builder{}

	la.serialize(&builder)
	return builder.String()
}

// Serialize returns the variable as a De Bruijn index representation
func (lv Var) Serialize() string {
	idx := int(lv) + 1
	return strconv.Itoa(idx)
}

// Serialize the free varaible as ' (prime) followed by the variable name
func (lf Free) Serialize() string {
	return "'" + string(lf)
}

func (lx Appl) serialize(builder *strings.Builder) {
	// This flag is ued later on to decide if a space is necessary
	bracketInMiddle := false

	switch lx[0].(type) {
	case Abst:
		builder.WriteByte('(')
		lx[0].serialize(builder)
		builder.WriteByte(')')

		bracketInMiddle = true

	default:
		lx[0].serialize(builder)
	}

	switch lx[1].(type) {
	case Appl:
		builder.WriteByte('(')
		lx[1].serialize(builder)
		builder.WriteByte(')')

	default:
		if !bracketInMiddle {
			builder.WriteByte(' ')
		}

		lx[1].serialize(builder)
	}
}

func (la Abst) serialize(builder *strings.Builder) {
	builder.WriteByte('l')
	la[0].serialize(builder)
}

func (lv Var) serialize(builder *strings.Builder) {
	idx := int(lv) + 1
	builder.WriteString(strconv.Itoa(idx))
}

// serialize the free varaible as ' (prime) followed by the variable name
func (lf Free) serialize(builder *strings.Builder) {
	builder.WriteByte('\'')
	builder.WriteString(string(lf))
}
