package post

import (
	"time"
	"github.com/VictoriaOtm/forum-api/models/user"
	"github.com/VictoriaOtm/forum-api/models/forum"
	"github.com/VictoriaOtm/forum-api/models/thread"
)

type Post struct {
	Author   string
	Created  time.Time
	Forum    string
	Id       int64
	IsEdited bool `json:"isEdited,omitempty"`
	Parent   int64
	Thread   int32
	Message  string
}

type PostDetails struct {
	Post   Post
	User   *user.User     `json:"author,omitempty"`
	Forum  *forum.Forum   `json:",omitempty"`
	Thread *thread.Thread `json:",omitempty"`
}

type PostUpdate struct {
	Message *string
}

//easyjson:json
type PostsArr []Post
