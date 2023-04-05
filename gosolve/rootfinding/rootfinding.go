package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/blldr/freshgoutils/stack"
	"github.com/blldr/gosolve/gosolve/mathparse"
)

func FindRoots(text string, iStart float64, iEnd float64, resolution float64, epsion float64)  ([]float64, error) {
	eqs := strings.Split(text, "=") 
	eqLen := len(eqs)
	if eqLen > 2 {
		return []float64{}, TooMuchEqualsSign{}
	}
	if eqLen == 1 {
		fmt.Println("first")
		eq, err := mathparse.ParseExpression(text)
		if err != nil {
			return nil, TooMuchEqualsSign{}
		}
		roots, err := findRoots(*eq, iStart, iEnd, resolution, epsion)
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
	eq2Len := eq2.Length()
	eq2Inverse := stack.NewStack[string]()
	for i := 0; i < eq2Len; i++ {
		eq2Inverse.Push(*eq2.Pop())
	}
	for i := 0; i < eq2Len; i++ {
		eq1.Push(*eq2Inverse.Pop())
	}
	eq1.Push("-")
	roots, err := findRoots(*eq1, iStart, iEnd, resolution, epsion) 
	if err != nil {
		return nil, err
	}
	return roots, nil
	
}

func findRoots(eq stack.Stack[string], iStart float64, iEnd float64, resolution float64, epsion float64) ([]float64, error) {
	a, err := mathparse.EvalExpression(eq, map[rune]float64{'y':iStart})
	if err != nil {
		return nil, err
	}
	b, err := mathparse.EvalExpression(eq, map[rune]float64{'y':iStart + resolution})
	if err != nil {
		return nil, err
	}
	roots := make([]float64, 0, 10)
	for i := 1.0;iStart + i*resolution < iEnd; i++ {
		fmt.Println(i)
		if (math.Signbit(a) == math.Signbit(b)) {
			fmt.Println(iStart + i*resolution)
			a, err = mathparse.EvalExpression(eq, map[rune]float64{'y':iStart + i*resolution})
			if err != nil {
				return nil, err
			}
			b, err = mathparse.EvalExpression(eq, map[rune]float64{'y':iStart + (i+1)*resolution})
			if err != nil {
				return nil, err
			}
			continue
		}
		c,d := a,b
		cX, dX := iStart + i*resolution, iStart + (i+1)*resolution
		fmt.Println("DDD", c,d)
		return nil, nil
		iRes := resolution / 2
		for (math.Abs(c-d) > epsion) {
			fmt.Println(c, d)
			if (math.Abs(c) > math.Abs(d)) {
				cX = cX + iRes
				c, _ = mathparse.EvalExpression(eq, map[rune]float64{'y': cX})
			} else {
				dX = dX - iRes
				d, _ = mathparse.EvalExpression(eq, map[rune]float64{'y': dX})
				
			}
			iRes = iRes / 2
			
		}
		roots = append(roots, (c + d)/2)
	}
	return roots, nil
}
