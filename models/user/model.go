package user

type User struct {
	About    string
	Email    string
	Fullname string
	Nickname string
}

type Update struct {
	About    *string
	Email    *string
	Fullname *string
}

//easyjson:json
type Arr []User
