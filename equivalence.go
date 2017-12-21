package LamCalc

// Currently someterm.Equals(otherterm) only uses alphaEquivalence
// should this include other checks?

// Equivalent checks if the LamExpr and LamTerm are Equivalent
func (lx LamExpr) Equivalent(other LamTerm) bool {
	return lx.alphaEquivalent(other)
}

// Equivalent checks if the LamAbst and LamTerm are Equivalent
func (la LamAbst) Equivalent(other LamTerm) bool {
	return la.alphaEquivalent(other)
}

// Equivalent checks if the LamVar and LamTerm are Equivalent
func (lv LamVar) Equivalent(other LamTerm) bool {
	return lv.alphaEquivalent(other)
}
