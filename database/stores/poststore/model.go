package poststore

import (
	"time"
)

//go:generate easyjson -snake_case

//easyjson:json
type Post struct {
	Author   string
	Created  time.Time
	Forum    string
	Id       int64
	IsEdited bool `json:"isEdited,omitempty"`
	Parent   int64
	Thread   int32
	Message  string
	parents  []int64
}

func (p *Post) MustUnmarshalJSON(b []byte) {
	err := p.UnmarshalJSON(b)
	if err != nil {
		panic(err)
	}
}

func (p *Post) MustMarshalJSON() []byte {
	b, err := p.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

//easyjson:json
type PostUpdate struct {
	Message *string
}

func (p *PostUpdate) MustUnmarshalJSON(b []byte) {
	err := p.UnmarshalJSON(b)
	if err != nil {
		panic(err)
	}
}

func (p *PostUpdate) MustMarshalJSON() []byte {
	b, err := p.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

//easyjson:json
type PostSlice []Post

func (p *PostSlice) MustUnmarshalJSON(b []byte) {
	err := p.UnmarshalJSON(b)
	if err != nil {
		panic(err)
	}
}

func (p *PostSlice) MustMarshalJSON() []byte {
	b, err := p.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}
