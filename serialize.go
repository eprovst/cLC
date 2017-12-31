package LamCalc

import (
	"bytes"
	"strconv"
)

// Serialize returns the application as a De Bruijn index representation
func (lx Appl) Serialize() string {
	buffer := bytes.NewBufferString("")

	lx.serialize(buffer)
	return buffer.String()
}

// Serialize returns the lambda abstraction as a De Bruijn index representation
func (la Abst) Serialize() string {
	buffer := bytes.NewBufferString("")

	la.serialize(buffer)
	return buffer.String()
}

// Serialize returns the variable as a De Bruijn index representation
func (lv Var) Serialize() string {
	idx := int(lv) + 1
	return strconv.Itoa(idx)
}

func (lx Appl) serialize(buffer *bytes.Buffer) {
	// This flag is ued later on to decide if a space is necessary
	bracketInMiddle := false

	switch lx[0].(type) {
	case Abst:
		buffer.WriteByte('(')
		lx[0].serialize(buffer)
		buffer.WriteByte(')')

		bracketInMiddle = true

	default:
		lx[0].serialize(buffer)
	}

	switch lx[1].(type) {
	case Appl:
		buffer.WriteByte('(')
		lx[1].serialize(buffer)
		buffer.WriteByte(')')

	default:
		if !bracketInMiddle {
			buffer.WriteByte(' ')
		}

		lx[1].serialize(buffer)
	}
}

func (la Abst) serialize(buffer *bytes.Buffer) {
	buffer.WriteByte('l')
	la[0].serialize(buffer)
}

func (lv Var) serialize(buffer *bytes.Buffer) {
	idx := int(lv) + 1
	buffer.WriteString(strconv.Itoa(idx))
}
