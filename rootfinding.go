package gosolve

import (
	"math"
	"strings"

	"github.com/blldr/freshgoutils/stack"
)

func FindRoots(text string, iStart float64, iEnd float64, variableName rune) ([]float64, error) {
	eqs := strings.Split(text, "=")
	eqLen := len(eqs)
	if eqLen > 2 {
		return nil, TooMuchEqualsSign{}
	}
	if eqLen == 1 {
		eq, err := ParseExpression(text)
		if err != nil {

			return nil, TooMuchEqualsSign{}
		}
		roots, err := findRoots(*eq, iStart, iEnd, variableName)
		if err != nil {
			return nil, err
		}
		return roots, nil
	}
	eq1, err := ParseExpression(eqs[0])
	if err != nil {
		return nil, err
	}
	eq2, err := ParseExpression(eqs[1])
	if err != nil {
		return nil, err
	}
	tmpStack := stack.NewStack[string]()
	eq1Len := eq1.Length()
	eq2Len := eq2.Length()
	for i := 0; i < eq1Len; i++ {
		tmpStack.Push(*eq1.Pop())
	}
	for i := 0; i < eq2Len; i++ {
		tmpStack.Push(*eq2.Pop())
	}
	eqStack := stack.NewStack[string]()
	eqStack.Push("-")
	tmpLen := tmpStack.Length()
	for i := 0; i < tmpLen; i++ {
		eqStack.Push(*tmpStack.Pop())
	}
	roots, err := findRoots(*eqStack, iStart, iEnd, variableName)
	if err != nil {
		return nil, err
	}
	return roots, nil

}

func findRoots(eq stack.Stack[string], iStart float64, iEnd float64, variableName rune) ([]float64, error) {
	epsilon := 0.01
	resolution := 0.1
	roots := make([]float64, 0, 10)
	a, b := iStart, iStart+resolution
	for a <= iEnd {
		if a > b {
			a, b = b, a
		}
		fa, err := EvalExpression(eq, map[rune]float64{variableName: a})
		if err != nil {
			return nil, err
		}
		fb, err := EvalExpression(eq, map[rune]float64{variableName: b})

		if err != nil {
			return nil, err
		}
		if fb == 0 {
			roots = append(roots, b)
			a = b
			b = a + resolution
			continue
		}
		if fa*fb < 0 {
			absFA := math.Abs(fa)
			absFB := math.Abs(fb)
			interval := b - a
			var tan float64
			tan = (absFB + absFA) / (b - a)
			if fa > fb {
				tan = -tan
			}
			var d float64
			if tan > 0 {
				d = b - absFB/(absFA+absFB)*interval
			} else {
				d = a + absFA/(absFA+absFB)*interval
			}
			d = roundFloat(d, 3)
			fd, err := EvalExpression(eq, map[rune]float64{variableName: d})
			if err != nil {
				return nil, err
			}
			if math.Abs(fd-0) < epsilon {
				roots = append(roots, roundFloat(d, 2))
				a = b
				b = a + resolution
				continue
			}

			if d == a {
				if !(math.Abs(0-tan) > 9000) {
					roots = append(roots, roundFloat(d, 2))
				}
				a = b
				b = a + resolution
				continue
			}
			if d == b {
				if !(math.Abs(0-tan) > 9000) {
					roots = append(roots, roundFloat(d, 2))
				}
				a = b
				b = a + resolution
				continue
			}

			if fd > 0 {
				if tan > 0 {
					b = d
					continue
				}
				a = d
			} else {
				if tan > 0 {
					a = d
					continue
				}
				b = d
			}
			continue
		}
		a = b
		b = a + resolution
	}

	return roots, nil
}

func roundFloat(f float64, precission int) float64 {
	return math.Round(f*math.Pow10(precission)) / math.Pow10(precission)
}

func floorFloat(f float64, precission int) float64 {
	return math.Floor(f*math.Pow10(precission)) / math.Pow10(precission)
}
