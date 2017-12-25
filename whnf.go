package LamCalc

// WHNF encapsulates the expression in a lambda abstraction
func (lx LamExpr) WHNF() LamAbst {
	nw := lx.Simplify()

	switch nw := nw.(type) {
	case LamExpr:
		// The simplification of a LamAbst is always a LamAbst
		return LamAbst{heightenIndex(lx), LamVar(0)}.Simplify().(LamAbst)

	default:
		return nw.WHNF()
	}
}

// WHNF returns the abstraction
func (la LamAbst) WHNF() LamAbst {
	// The simplification of a LamAbst is always a LamAbst
	return la.Simplify().(LamAbst)
}

// WHNF encapsulates the Lambda variable inside of a lambda abstraction
func (lv LamVar) WHNF() LamAbst {
	return LamAbst{lv, LamVar(0)}
}
