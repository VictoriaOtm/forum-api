package model_pool

import (
	"sync"
	"github.com/VictoriaOtm/forum-api/models/post"
)

var PostArrPool = sync.Pool{
	New: func() interface{} {
		return make(post.PostsArr, 0, 300)
	},
}

var PostDetailsPool = sync.Pool{
	New: func() interface{} {
		return &post.PostDetails{}
	},
}
