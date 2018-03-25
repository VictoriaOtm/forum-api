package model_pool

import (
	"sync"
	"github.com/VictoriaOtm/forum-api/models/forum"
)

var ForumPool = sync.Pool{
	New: func() interface{} {
		return &forum.Forum{}
	},
}
