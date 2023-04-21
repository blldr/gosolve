package gosolve

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
		text  string
		index int
	}

	token struct {
		tokenType tokenType
		data      string
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

func (t *tokenizer) nextToken() (*token, error) {
	c := t.getChar()

	if c == 0 {
		return nil, EndOfString{}
	}

	if isDigit(c) {
		t.putbackChar()
		tokenData, err := t.parseNumeric()
		if err != nil {
			return nil, err
		}
		return &token{Numeric, tokenData}, nil
	}

	switch c {
	case '+':
		return &token{Operator, string(c)}, nil
	case '-':
		return &token{Operator, string(c)}, nil
	case '*':
		return &token{Operator, string(c)}, nil
	case '/':
		return &token{Operator, string(c)}, nil
	case '^':
		return &token{Operator, string(c)}, nil
	case '(':
		return &token{Operator, string(c)}, nil
	case ')':
		return &token{Operator, string(c)}, nil
	}

	if isLetter(c) {
		t.putbackChar()
		tokenData := t.parseFunction()
		if len(tokenData) == 1 {
			return &token{Variable, tokenData}, nil
		}
		return &token{Function, tokenData}, nil
	}

	return nil, InvalidToken{}
}

func (t *tokenizer) parseNumeric() (string, error) {
	var b strings.Builder
	c := t.getChar()

	for i := 0; isDigit(c); i++ {
		b.WriteRune(c)
		c = t.getChar()
	}

	if c == '.' {
		b.WriteRune('.')
		c = t.getChar()
		if isDigit(c) {
			for i := 0; isDigit(c); i++ {
				b.WriteRune(c)
				c = t.getChar()
			}
		} else {
			return "", InvalidToken{}
		}
	}

	if c != 0 {
		t.putbackChar()
	}

	return string(b.String()), nil
}

func (t *tokenizer) parseFunction() string {
	var b strings.Builder

	c := t.getChar()
	for i := 0; unicode.IsLetter(c); i++ {
		b.WriteRune(c)
		c = t.getChar()
	}

	if c != 0 {
		t.putbackChar()
	}

	return b.String()
}

func isDigit(c rune) bool {
	if c >= 48 && c <= 57 {
		return true
	}

	return false
}

func isLetter(c rune) bool {
	if (c >= 65 && c <= 90) || (c >= 97 && c <= 122) {
		return true
	}

	return false
}
