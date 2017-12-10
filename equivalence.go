package LamCalc

// Currently someterm.Equals(otherterm) only uses alphaEquivalence
// should this include other checks?

// Equivalent checks if the LamExpr and LamTerm are Equivalent
func (lx LamExpr) Equivalent(other LamTerm) bool {
	return lx.alphaEquivalent(other)
}

// Equivalent checks if the LamAbst and LamTerm are Equivalent
func (lf LamAbst) Equivalent(other LamTerm) bool {
	return lf.alphaEquivalent(other)
}
