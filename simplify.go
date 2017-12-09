package LamCalc

import "reflect"

// Simplify (tries) to remove unnecessary brackets
func (lx LamExpr) simplify() LamTerm {
	if len(lx) == 1 {
		if reflect.TypeOf(lx[0]).Kind() != reflect.Int {
			return lx[0].(LamTerm).simplify()
		}

		return lx

	} else if reflect.TypeOf(lx[0]).String() == "LamCalc.LamExpr" {
		res := lx[0].(LamExpr)

		if len(lx) > 1 {
			res = append(res, lx[1:]...)
		}

		return res.simplify()
	}

	res := LamExpr{}

	for _, term := range lx {
		switch term := term.(type) {
		case LamTerm:
			switch simpl := term.simplify().(type) {
			case LamExpr:
				if len(simpl) == 1 {
					res = append(res, simpl[0])
				} else {
					res = append(res, simpl)
				}

			case LamFunc:
				res = append(res, simpl)
			}

		default:
			res = append(res, term)
		}
	}

	return res
}

// Simplify (tries) to remove unnecessary brackets
func (lf LamFunc) simplify() LamTerm {
	simpl := LamExpr(lf).simplify()

	switch simpl := simpl.(type) {
	case LamExpr:
		return LamFunc(simpl)

	default:
		return LamFunc{simpl.(LamFunc)}
	}
}
