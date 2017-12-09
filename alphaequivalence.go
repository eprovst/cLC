package LamCalc

import "reflect"

// As we are using De Bruijn indexes alphaEquivalence is the same as syntactic equivalence

// Equals checks wether the LamExp is identical to a LamTerm
func (lx LamExpr) alphaEquivalent(other LamTerm) bool {

	if reflect.TypeOf(other).String() != "LamCalc.LamExpr" {
		return false

	} else if len(lx) != other.len() {
		return false
	}

	for i := range lx {
		switch elem := lx[i].(type) {
		case int:
			if reflect.TypeOf(other.index(i)).Kind() != reflect.Int || elem != other.index(i).(int) {
				return false
			}

		case LamExpr:
			if reflect.TypeOf(other.index(i)).String() != "LamCalc.LamExpr" || !elem.alphaEquivalent(other.index(i).(LamExpr)) {
				return false
			}

		case LamFunc:
			if reflect.TypeOf(other.index(i)).String() != "LamCalc.LamFunc" || !LamExpr(elem).alphaEquivalent(LamExpr(other.index(i).(LamFunc))) {
				return false
			}

		default:
			return false
		}
	}

	return true
}

// Equals checks wether a LamFunc and a LamTerm are identical
func (lf LamFunc) alphaEquivalent(other LamTerm) bool {

	if reflect.TypeOf(other).String() != "LamCalc.LamFunc" {
		return false
	}

	return LamExpr(lf).alphaEquivalent(LamExpr(other.(LamFunc)))
}
