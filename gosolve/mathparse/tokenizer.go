package mathparse

import (
	"strings"
	"unicode"
)

type tokenType int
const (
	Numeric tokenType = iota
	Operator
	Function
	Variable
)

type (
	tokenizer struct {
		text string
		index int
	}
) 

func newTokenizer(text string) *tokenizer {
	text = strings.ReplaceAll(text, " ", "")
	return &tokenizer{
		text, 
		0,
	}
}

func (t *tokenizer) getChar() rune {
	if t.index >= len(t.text) {
		return 0
	}

	tmp := t.index
	t.index += 1
	return rune(t.text[tmp])

}

func (t *tokenizer) putbackChar() {
	t.index -= 1
}

func (t *tokenizer) nextToken() (string, tokenType, error) {
	c := t.getChar()
	if c == 0 {
		return "", 0, EndOfString{}
	}
	switch true {
	case unicode.IsDigit(c):
		t.putbackChar()
		num := t.parseNumeric()
		return num, Numeric, nil
	case c == '+':
		return string(c), Operator, nil
	case c == '-':
		return string(c), Operator, nil
	case c == '*':
		return string(c), Operator, nil
	case c == '/':
		return string(c), Operator, nil
	case c == '^':
		return string(c), Operator, nil
	case c == '(':
		return string(c), Operator, nil
	case c == ')':
		return string(c), Operator, nil
	case unicode.IsLetter(c):
		t.putbackChar()
		fun := t.parseFunction()
		if len(fun) == 1 {
			return fun, Variable, nil
		}
		return fun, Function, nil
	}
	return "", 0, InvalidToken{}
}

func (t *tokenizer) parseNumeric() string {
	buffer := make([]rune, 0, 20)
	c := t.getChar()
	for i := 0; unicode.IsDigit(c); i++ {
		buffer = append(buffer, c)
		c = t.getChar()
	}
	if c == '.' {
		buffer = append(buffer, '.')
		c = t.getChar()
		if unicode.IsNumber(c) {
			for i := 0; unicode.IsDigit(c); i++ {
				buffer = append(buffer, c)
				c = t.getChar()
			}
			
		}
	}
	if c != 0 {
		t.putbackChar()
	}
	return string(buffer)
}

func (t *tokenizer) parseFunction() string{
	buffer := make([]rune, 0, 20)
	c := t.getChar()
	for i := 0; unicode.IsLetter(c); i++ {
		buffer = append(buffer, c)
		c = t.getChar()
	}
	if c != 0 {
		t.putbackChar()
	}
	return string(buffer)
}
