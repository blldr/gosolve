package mathparse

import (
	"fmt"
	"math"
	"strings"

	"github.com/blldr/freshgoutils/stack"
)

func FindRoots(text string, iStart float64, iEnd float64, variableName rune)  ([]float64, error) {
	eqs := strings.Split(text, "=") 
	eqLen := len(eqs)
	if eqLen > 2 {
		return []float64{}, TooMuchEqualsSign{}
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
	for i := 0; i < eq1Len; i++{
		tmpStack.Push(*eq1.Pop())
	}
	for i := 0; i < eq2Len; i++{
		tmpStack.Push(*eq2.Pop())
	}
	eqStack := stack.NewStack[string]()
	eqStack.Push("-")
	tmpLen := tmpStack.Length()
	for i := 0; i < tmpLen; i++ {
		eqStack.Push(*tmpStack.Pop())
	}
	roots, err := findRoots(*eqStack, iStart, iEnd, 'y')
	return roots, nil
	
}

func findRoots(eq stack.Stack[string], iStart float64, iEnd float64, variableName rune) ([]float64, error) {
	epsilon := 0.1
	resolution := 0.2
	roots := make([]float64, 0, 10)
	a, b := iStart, iStart + resolution
	for ; a <= iEnd;  {
		a = roundFloat(a, 2)
		b = roundFloat(b, 2)
		fa, err := EvalExpression(eq, map[rune]float64{variableName:a})
		if err != nil {
			return nil, err
		}
		fa = roundFloat(fa, 2)
		fb, err := EvalExpression(eq, map[rune]float64{variableName:b})
		if err != nil {
			return nil, err
		}
		fb = roundFloat(fb, 2)
		if fb == 0 {
			roots = append(roots, b)
			a = b 
			b = a + resolution
			continue
		}

		if fa * fb < 0 {
			absFA := math.Abs(fa)
			absFB := math.Abs(fb)
			interval := b - a
			var tan bool
			if fa < 0 {
				tan = true
			} else {
				tan = false
			}
			d := roundFloat(b - absFB / (absFA + absFB) * interval, 2)
			fd, err := EvalExpression(eq, map[rune]float64{variableName: d})
			if err != nil {
				return nil, err
			}
			fd = roundFloat(fd, 2)
			if math.Abs(fd - 0) < epsilon {
				roots = append(roots, roundFloat(d, 2))
				a = b
				b = a + resolution
				continue
			}
			fmt.Println(a, b, d)
			if d == a {
				roots = append(roots, roundFloat(a,2))
				a = b
				b = a + resolution
				continue
			}
			if d == b {
				roots = append(roots, roundFloat( b - epsilon,2))
				a = b
				b = a + resolution
			}

			if fd > 0 {
				if tan {
					b = d
					continue
				}
				a = d
			} else {
				if tan {
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

func roundFloat(f float64, precission int) float64{
	return math.Round(f * math.Pow10(precission)) / math.Pow10(precission)
}
