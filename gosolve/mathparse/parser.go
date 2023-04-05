package mathparse

import (
	"math"
	"strconv"
	"unicode"

	"github.com/blldr/freshgoutils/stack"
)

type OperationType int
const (
	Add OperationType = iota + 1
	Sub
	Mul
	Div
	Sin
	Cos
)

type RelativeExpression int 
const (
	Higher RelativeExpression = iota
	Same
	Lower
)



func ParseExpression(text string) (*stack.Stack[string], error) {
	operatorStack := stack.NewStack[string]()
	outputStack := stack.NewStack[string]() 
	tokenizer := newTokenizer(text)
	for token, tokenType, err := tokenizer.nextToken(); err == nil; token, tokenType, err = tokenizer.nextToken() {
		switch tokenType {
		case Numeric:
		outputStack.Push(token)
		case Variable:
		outputStack.Push(token)
		case Operator:
			if token == "(" {
				operatorStack.Push("(")
				continue
			}
			if token == ")" {
				o := operatorStack.Pop()
				if o == nil {
					break
				}
				for(*o != "(") {
					
					outputStack.Push(*o)
					o = operatorStack.Pop()
					if o == nil {
						break
					}
				}
				break
			}
			for(true) {
				a := operatorStack.Pop()
				if a == nil {
					break
				} 
				if *a == "(" {
					operatorStack.Push("(")
					break
				}
				cmpResult := cmpOp(token, *a)
				if cmpResult == Higher {
					outputStack.Push(*a)
					continue
				}
				if cmpResult == Same {
					outputStack.Push(*a)
					continue
				}
				operatorStack.Push(*a)
				break
			}
			operatorStack.Push(token)
			continue
		case Function:
			if isFunc(token) {
				operatorStack.Push(token)
				continue
			}
			return nil, UndefinedFunction{}
		}	

	}
	len := operatorStack.Length()
	for i := 0; i < len; i ++ {
		outputStack.Push(*operatorStack.Pop())
	}
	inverseStack := stack.NewStack[string]()
	len = outputStack.Length()
	for i := 0; i < len; i ++ {

		inverseStack.Push(*outputStack.Pop())
		
	}
	return inverseStack, nil
}

func EvalExpression(eqStack stack.Stack[string], vars map[rune]float64) (float64, error){
	numbersStack := stack.NewStack[float64]()
	slen := eqStack.Length()
	for i := 0; i < slen; i ++ {
		a := *eqStack.Pop()
		if unicode.IsNumber(rune(a[0])) {
			f, err := strconv.ParseFloat(a, 64)
			if err != nil {
				return 0, UnexpectedError{}
			}
			numbersStack.Push(f)
			continue
		}

		if unicode.IsLetter(rune(a[0])) && len(a)  == 1 {
			val, ok := vars[rune(a[0])]
			if !ok {
				return 0, UndefinedVariable{}
			}
			numbersStack.Push(float64(val))
			continue
		}
		op1 := numbersStack.Pop()
		if op1 == nil {
			return 0, UnexpectedError{}
		}
		if isFunc(a) {
			switch a {
			case "sin":
				numbersStack.Push(math.Sin(*op1))
			case "cos":
				numbersStack.Push(math.Cos(*op1))
			}
			continue
		}
		op2 := numbersStack.Pop()
		if op2 == nil {
			return 0, UnexpectedError{}
		}
		switch a {
		case "+":
			numbersStack.Push(*op2 + *op1)
		case "-":
			numbersStack.Push(*op2 - *op1)
		case "*":
			numbersStack.Push(*op2 * *op1)
		case "/":
			numbersStack.Push(*op2 / *op1)
		case "^":
		numbersStack.Push(math.Pow(*op2, *op1))

		}
		continue
	}
	if numbersStack.Length() != 1 {
		return 0, UnexpectedError{}
	}
	return *numbersStack.Pop(), nil
}


func isFunc(op string) bool {
	switch op {
	case "sin":
		return true
	case "cos":
		return true
	}
	return false
}

func cmpOp(op1 string, op2 string) RelativeExpression {
	if op1 == "+" || op1 == "-" {
		if op2 == "*" || op2 == "/" || op2 == "^" {
			return Higher
		}
		return Same
	}
	if op1 == "*" || op1 == "/" {
		if op2 == "^" {
			return Higher
		}
		if op2 == "*" || op2 == "/" {
			return Same
		}
		return Lower
	}
	if op1 == "^" && op2 == "^" {
		return Same
	}
	return Lower

}

