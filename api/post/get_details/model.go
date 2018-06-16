package get_details

//go:generate easyjson -snake_case

import (
	"sync"

	"github.com/VictoriaOtm/forum-api/database/stores/forumstore"
	"github.com/VictoriaOtm/forum-api/database/stores/poststore"
	"github.com/VictoriaOtm/forum-api/database/stores/threadstore"
	"github.com/VictoriaOtm/forum-api/database/stores/userstore"
)

//easyjson:json
type postDetails struct {
	Post   poststore.Post
	User   *userstore.User     `json:"author,omitempty"`
	Forum  *forumstore.Forum   `json:",omitempty"`
	Thread *threadstore.Thread `json:",omitempty"`
}

func (pd *postDetails) MustMarshalJSON() []byte {
	b, err := pd.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

type postDetailsPool struct {
	sync.Pool
}

func (p *postDetailsPool) Acquire() *postDetails {
	return p.Get().(*postDetails)
}

func (p *postDetailsPool) Utilize(d *postDetails) {
	d.User = nil
	d.Forum = nil
	d.Thread = nil

	p.Put(d)
}

var pdPool = postDetailsPool{
	sync.Pool{
		New: func() interface{} {
			return &postDetails{}
		},
	},
}
