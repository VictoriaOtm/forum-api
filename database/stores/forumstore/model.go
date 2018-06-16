package forumstore

//go:generate easyjson -snake_case

//Forum is model representing table in db
//easyjson:json
type Forum struct {
	Slug    string
	Title   string
	User    string
	Posts   int32 `json:",omitempty"`
	Threads int32 `json:",omitempty"`
}

func (f *Forum) MustUnmarshalJSON(bs []byte) {
	err := f.UnmarshalJSON(bs)
	if err != nil {
		panic(err)
	}
}

func (f *Forum) MustMarshalJSON() []byte {
	bs, err := f.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return bs
}
