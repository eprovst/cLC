package LamCalc

// substitute replaces index by sub
func (lx LamExpr) substitute(index LamVar, sub LamTerm) LamTerm {
	nw := LamExpr{}

	for _, term := range lx {
		switch term := term.(type) {
		case LamAbst:
			nw = append(nw, term.substitute(index+1, heightenIndex(sub)))

		default:
			nw = append(nw, term.substitute(index, sub))
		}
	}

	return nw
}

// substitute replaces index by sub
func (la LamAbst) substitute(index LamVar, sub LamTerm) LamTerm {
	nw := LamAbst{}

	for _, term := range la {
		switch term := term.(type) {
		case LamAbst:
			nw = append(nw, term.substitute(index+1, heightenIndex(sub)))

		default:
			nw = append(nw, term.substitute(index, sub))
		}
	}

	return nw
}

// substitute replaces index by sub
func (lv LamVar) substitute(index LamVar, sub LamTerm) LamTerm {
	if lv == index {
		return sub
	}

	return lv
}

// betaReduce replaces index 0 by sub and returns a LamExpr
func (la LamAbst) betaReduce(sub LamTerm) LamTerm {
	return lowerIndex(
		LamExpr(la).substitute(0, heightenIndex(sub)),
	)
}
