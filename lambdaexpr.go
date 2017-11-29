package LamCalc

import "reflect"

// LamExpr is a list of lamfuncs, lamexprns and De Bruijn indexes (all lowered by one) which isn't a function itself.
type LamExpr []interface{}

// Equals checks wether two LamExprns are identical.
func (lx LamExpr) Equals(other LamExpr) bool {
	if len(lx) != len(other) {
		return false
	}

	for i := range lx {
		switch elem := lx[i].(type) {
		case int:
			if reflect.TypeOf(other[i]).Kind() != reflect.Int || elem != other[i].(int) {
				return false
			}

		case LamExpr:
			if reflect.TypeOf(other[i]).String() != "LamCalc.LamExpr" || !elem.Equals(other[i].(LamExpr)) {
				return false
			}

		case LamFunc:
			if reflect.TypeOf(other[i]).String() != "LamCalc.LamFunc" || !LamExpr(elem).Equals(LamExpr(other[i].(LamFunc))) {
				return false
			}

		default:
			return false
		}
	}

	return true
}

// LamFunc is a list of lamfuncs, lamexprns and De Bruijn indexes (all lowered by one).
type LamFunc []interface{}

// Equals checks wether two LamFuncs are identical.
func (lf LamFunc) Equals(other LamFunc) bool {
	return LamExpr(lf).Equals(LamExpr(other))
}
