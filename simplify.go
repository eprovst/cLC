package LamCalc

import "reflect"

// Simplify (tries) to remove unnecessary brackets
func (lx LamExpr) Simplify() LamTerm {
	if len(lx) == 1 {
		return lx[0].Simplify()

	} else if reflect.TypeOf(lx[0]).String() == "LamCalc.LamExpr" {
		res := lx[0].(LamExpr)

		if len(lx) > 1 {
			res = append(res, lx[1:]...)
		}

		return res.Simplify()
	}

	res := LamExpr{}

	for _, term := range lx {
		switch term := term.Simplify().(type) {
		case LamExpr:
			if len(term) == 1 {
				res = append(res, term[0])
			} else {
				res = append(res, term)
			}

		default:
			res = append(res, term)
		}
	}

	return res
}

// Simplify (tries) to remove unnecessary brackets
func (la LamAbst) Simplify() LamTerm {
	simpl := LamExpr(la).Simplify()

	switch simpl := simpl.(type) {
	case LamExpr:
		return LamAbst(simpl)

	default:
		return LamAbst{simpl}
	}
}

// Simplify (tries) to remove unnecessary brackets
func (lv LamVar) Simplify() LamTerm {
	return lv
}
