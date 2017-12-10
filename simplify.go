package LamCalc

import "reflect"

// Simplify (tries) to remove unnecessary brackets
func (lx LamExpr) Simplify() LamTerm {
	if len(lx) == 1 {
		if reflect.TypeOf(lx[0]).Kind() != reflect.Int {
			return lx[0].(LamTerm).Simplify()
		}

		return lx

	} else if reflect.TypeOf(lx[0]).String() == "LamCalc.LamExpr" {
		res := lx[0].(LamExpr)

		if len(lx) > 1 {
			res = append(res, lx[1:]...)
		}

		return res.Simplify()
	}

	res := LamExpr{}

	for _, term := range lx {
		switch term := term.(type) {
		case LamTerm:
			switch simpl := term.Simplify().(type) {
			case LamExpr:
				if len(simpl) == 1 {
					res = append(res, simpl[0])
				} else {
					res = append(res, simpl)
				}

			case LamAbst:
				res = append(res, simpl)
			}

		default:
			res = append(res, term)
		}
	}

	return res
}

// Simplify (tries) to remove unnecessary brackets
func (lf LamAbst) Simplify() LamTerm {
	simpl := LamExpr(lf).Simplify()

	switch simpl := simpl.(type) {
	case LamExpr:
		return LamAbst(simpl)

	default:
		return LamAbst{simpl.(LamAbst)}
	}
}
