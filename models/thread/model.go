package thread

import "time"

type Thread struct {
	Author  string
	Created time.Time
	Forum   string
	Id      int32
	Message string
	Slug    *string `json:"slug,omitempty"`
	Title   string
	Votes   int32   `json:"votes,omitempty"`
}

type Vote struct {
	Nickname string
	Voice    int32
}

type Update struct {
	Message *string
	Title   *string
}

//easyjson:json
type Arr []Thread
