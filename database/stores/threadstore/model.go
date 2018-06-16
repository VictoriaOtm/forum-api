package threadstore

import (
	"time"
)

//go:generate easyjson -snake_case

//easyjson:json
type Thread struct {
	Author  string
	Created time.Time
	Forum   string
	Id      int32
	Message string
	Slug    *string `json:"slug,omitempty"`
	Title   string
	Votes   int32 `json:"votes,omitempty"`
}

func (t *Thread) MustUnmarshalJSON(b []byte) {
	err := t.UnmarshalJSON(b)
	if err != nil {
		panic(err)
	}
}

func (t *Thread) MustMarshalJSON() []byte {
	r, err := t.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return r
}

//easyjson:json
type ThreadSlice []Thread

func (ts *ThreadSlice) MustUnmarshalJSON(b []byte) {
	err := ts.UnmarshalJSON(b)
	if err != nil {
		panic(err)
	}
}

func (ts *ThreadSlice) MustMarshalJSON() []byte {
	b, err := ts.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

//easyjson:json
type ThreadUpdate struct {
	Message *string
	Title   *string
}

func (ts *ThreadUpdate) MustUnmarshalJSON(b []byte) {
	err := ts.UnmarshalJSON(b)
	if err != nil {
		panic(err)
	}
}

func (ts *ThreadUpdate) MustMarshalJSON() []byte {
	b, err := ts.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

//easyjson:json
type Vote struct {
	Nickname string
	Voice    int16
}
