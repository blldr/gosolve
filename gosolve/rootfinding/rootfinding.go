package main

import (
	"math"
	"strings"

	"github.com/blldr/freshgoutils/stack"
	"github.com/blldr/gosolve/gosolve/mathparse"
)

func FindRoots(text string, iStart float64, iEnd float64, variableName rune)  ([]float64, error) {
	eqs := strings.Split(text, "=") 
	eqLen := len(eqs)
	if eqLen > 2 {
		return []float64{}, TooMuchEqualsSign{}
	}
	if eqLen == 1 {
		eq, err := mathparse.ParseExpression(text)
		if err != nil {
			return nil, TooMuchEqualsSign{}
		}
		roots, err := findRoots(*eq, iStart, iEnd, variableName)
		if err != nil {
			return nil, err
		}
		return roots, nil
	}
	eq1, err := mathparse.ParseExpression(eqs[0])
	if err != nil {
		return nil, err
	}
	eq2, err := mathparse.ParseExpression(eqs[1])
	if err != nil {
		return nil, err
	}
	tmpStack := stack.NewStack[string]()
	eq1Len := eq1.Length()
	eq2Len := eq2.Length()
	for i := 0; i < eq2Len; i++{
		tmpStack.Push(*eq2.Pop())
	}
	for i := 0; i < eq1Len; i++{
		tmpStack.Push(*eq1.Pop())
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
	resolution := 0.5
	roots := make([]float64, 0, 10)
	a, b := iStart, iStart + resolution
	var tan bool
	var tano bool
	for ; a <= iEnd;  {
		if a > b {
			a,b = b,a
		}
		fa, err := mathparse.EvalExpression(eq, map[rune]float64{variableName:a})
		fb, err := mathparse.EvalExpression(eq, map[rune]float64{variableName: b})
		if err != nil {
			return nil, err
		}
		if (fb == 0) {
			roots = append(roots, b)

			a = b
			b = a + resolution
			tan = !tan
			continue
		}
		if fa < fb {
			tano = tan
			tan = true
		} else {
			tano = tan
			tan = false
		}
		if (fa * fb < 0) {
			interval := b - a
			absfA := math.Abs(a)
			absfB := math.Abs(b)
			d := b - math.Floor((absfB / (absfB+absfA) * interval) * 100) / 100
			fd, _ := mathparse.EvalExpression(eq, map[rune]float64{variableName: d})
			if math.Abs(fd - 0) < epsilon {
				roots = append(roots, math.Floor(d*10) / 10)

				for true {
					fa, _ := mathparse.EvalExpression(eq, map[rune]float64{variableName:a})
					if math.Abs(fa - 0) < epsilon {
						a = a + resolution
						continue
					}
					b = a + resolution
					break
				}
				b = a + resolution
				a = b + epsilon
				continue
			}
			if fd > 0 {
				b = d
			}
			a = d
			continue
		}
		if tano != tan {
		}
		a = b  
		b = a + resolution
	}


	return roots, nil
}
