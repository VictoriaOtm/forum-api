package userstore

//go:generate easyjson -snake_case

//easyjson:json
type User struct {
	About    string
	Email    string
	Fullname string
	Nickname string
}

func (u *User) MustUnmarshalJSON(b []byte) {
	err := u.UnmarshalJSON(b)
	if err != nil {
		panic(err)
	}
}

func (u *User) MustMarshalJSON() []byte {
	b, err := u.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

//easyjson:json
type UserSlice []User

func (u *UserSlice) MustUnmarshalJSON(b []byte) {
	err := u.UnmarshalJSON(b)
	if err != nil {
		panic(err)
	}
}

func (u *UserSlice) MustMarshalJSON() []byte {
	b, err := u.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

//easyjson:json
type UserUpdate struct {
	About    *string
	Email    *string
	Fullname *string
}

func (u *UserUpdate) MustUnmarshalJSON(b []byte) {
	err := u.UnmarshalJSON(b)
	if err != nil {
		panic(err)
	}
}

func (u *UserUpdate) MustMarshalJSON() []byte {
	b, err := u.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}
