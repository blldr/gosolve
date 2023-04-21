package gosolve

import (
	"errors"
	"math"
	"strconv"
	"unicode"

	"github.com/blldr/freshgoutils/stack"
)

type (
	RelativeExpression int
)

const (
	Higher RelativeExpression = iota
	Same
	Lower
)

func ParseExpression(text string) (*stack.Stack[string], error) {
	operatorStack := stack.NewStack[string]()
	outputStack := stack.NewStack[string]()
	tokenizer := newTokenizer(text)
	lastToken := Operator
	var errOut error
	for {
		token, err := tokenizer.nextToken()
		errOut = err
		if err != nil {
			break
		}

		switch token.tokenType {
		case Numeric:
			outputStack.Push(token.data)
			lastToken = Numeric
		case Variable:
			outputStack.Push(token.data)
			lastToken = Variable
		case Operator:
			if lastToken == Operator {
				if token.data == "+" || token.data == "-" {
					outputStack.Push("0")
				}
			}
			if token.data == ")" {
				o := operatorStack.Pop()
				if o == nil {
					break
				}
				for *o != "(" {
					outputStack.Push(*o)
					o = operatorStack.Pop()
					if o == nil {
						break
					}
				}
				break
			}

			lastToken = Operator
			if token.data == "(" {
				operatorStack.Push("(")
				continue
			}
			for true {
				a := operatorStack.Pop()
				if a == nil {
					break
				}
				if *a == "(" {
					operatorStack.Push("(")
					break
				}
				cmpResult := cmpOp(token.data, *a)
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
			operatorStack.Push(token.data)
			continue
		case Function:
			lastToken = Function
			if isFunc(token.data) {
				operatorStack.Push(token.data)
				continue
			}
			return nil, UndefinedFunction{}
		}

	}
	if errOut != nil && !errors.Is(errOut, EndOfString{}) {
		return nil, errOut
	}
	len := operatorStack.Length()
	for i := 0; i < len; i++ {
		outputStack.Push(*operatorStack.Pop())
	}
	return inverseStack(outputStack), nil
}

func EvalExpression(eqStack stack.Stack[string], vars map[rune]float64) (float64, error) {
	tmpStack := stack.NewStack[float64]()
	slen := eqStack.Length()
	for i := 0; i < slen; i++ {
		a := *eqStack.Pop()
		if unicode.IsNumber(rune(a[0])) {
			f, err := strconv.ParseFloat(a, 64)
			if err != nil {
				return 0, UnexpectedError{}
			}
			tmpStack.Push(f)
			continue
		}

		if unicode.IsLetter(rune(a[0])) && len(a) == 1 {
			val, ok := vars[rune(a[0])]
			if !ok {
				return 0, UndefinedVariable{}
			}
			tmpStack.Push(float64(val))
			continue
		}
		op1 := tmpStack.Pop()
		if op1 == nil {
			return 0, UnexpectedError{}
		}
		if isFunc(a) {
			switch a {
			case "sin":
				tmpStack.Push(math.Sin(*op1))
			case "cos":
				tmpStack.Push(math.Cos(*op1))
			case "tan":
				tmpStack.Push(math.Tan(*op1))
			case "cot":
				tmpStack.Push(1 / math.Tan(*op1))
			}
			continue
		}
		op2 := tmpStack.Pop()
		if op2 == nil {
			return 0, UnexpectedError{}
		}
		switch a {
		case "+":
			tmpStack.Push(*op2 + *op1)
		case "-":
			tmpStack.Push(*op2 - *op1)
		case "*":
			tmpStack.Push(*op2 * *op1)
		case "/":
			tmpStack.Push(*op2 / *op1)
		case "^":
			tmpStack.Push(math.Pow(*op2, *op1))

		}
		continue
	}
	if tmpStack.Length() != 1 {
		return 0, UnexpectedError{}
	}
	return *tmpStack.Pop(), nil
}

func isFunc(op string) bool {
	switch op {
	case "sin":
		return true
	case "cos":
		return true
	case "tan":
		return true
	case "cot":
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

func inverseStack[T any](s *stack.Stack[T]) *stack.Stack[T] {
	newStack := stack.NewStack[T]()
	stackLen := s.Length()
	for i := 0; i < stackLen; i++ {
		newStack.Push(*s.Pop())
	}

	return newStack
}
