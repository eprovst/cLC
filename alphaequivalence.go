package LamCalc

// As we are using De Bruijn indexes alphaEquivalence is the same as syntactic equivalence

// AlphaEquivalent checks whether the Appl is identical to a Term
func (lx Appl) AlphaEquivalent(other Term) bool {
	switch other := other.(type) {
	case Appl:
		return lx[0].AlphaEquivalent(other[0]) && lx[1].AlphaEquivalent(other[1])

	default:
		return false
	}
}

// AlphaEquivalent checks whether a Abst and a Term are identical
func (la Abst) AlphaEquivalent(other Term) bool {
	switch other := other.(type) {
	case Abst:
		return la[0].AlphaEquivalent(other[0])

	default:
		return false
	}
}

// AlphaEquivalent checks whether a Var and a Term are identical
func (lv Var) AlphaEquivalent(other Term) bool {
	switch other := other.(type) {
	case Var:
		return lv == other

	default:
		return false
	}
}
