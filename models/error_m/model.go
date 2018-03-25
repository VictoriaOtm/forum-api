package error_m

type Error struct {
	Message string `json:"message"`
}

var CommonError, _ = Error{"Error"}.MarshalJSON()