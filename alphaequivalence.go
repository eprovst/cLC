package LamCalc

// As we are using De Bruijn indexes alphaEquivalence is the same as syntactic equivalence

// Equals checks wether the Appl is identical to a Term
func (lx Appl) alphaEquivalent(other Term) bool {
	switch other := other.(type) {
	case Appl:
		return lx[0].alphaEquivalent(other[0]) && lx[1].alphaEquivalent(other[1])

	default:
		return false
	}
}

// Equals checks wether a Abst and a Term are identical
func (la Abst) alphaEquivalent(other Term) bool {
	switch other := other.(type) {
	case Abst:
		return la[0].alphaEquivalent(other[0])

	default:
		return false
	}
}

// Equals checks wether a Var and a Term are identical
func (lv Var) alphaEquivalent(other Term) bool {
	switch other := other.(type) {
	case Var:
		return lv == other

	default:
		return false
	}
}
