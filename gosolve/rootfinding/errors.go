package main


type (
	TooMuchEqualsSign struct {}
)

func (e TooMuchEqualsSign) Error() string {
	return "Too much equals sign"
}
