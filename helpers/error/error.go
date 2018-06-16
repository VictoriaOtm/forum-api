package error

//go:generate easyjson -snake_case

//easyjson:json
type Error struct {
	Message string
}

func (e *Error) MustMarshalJSON() []byte {
	b, err := e.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

func MakeError(msg string) []byte {
	e := Error{
		Message: msg,
	}
	return e.MustMarshalJSON()
}
