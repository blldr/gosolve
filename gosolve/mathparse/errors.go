package mathparse


type (
	EndOfString struct {}
	InvalidToken struct {}
	UndefinedFunction struct {}
	InvalidParanthesisStructure struct {}
	UnexpectedError struct {}
	UndefinedVariable struct {}
)

func (e EndOfString) Error() string {
	return "End of the string"
}

func (e InvalidToken) Error() string {
	return "Invalid token on equation string"
}

func (e UndefinedFunction) Error() string {
	return "Undefined Function"
}

func (e InvalidParanthesisStructure) Error() string {
	return "Invalid Paranthesis structure"
}

func (e UnexpectedError) Error() string {
	return "Unexpected error"
}

func (e UndefinedVariable) Error() string {
	return "Undefined Variable"	
}
