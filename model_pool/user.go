package model_pool

import (
	"sync"
	"github.com/VictoriaOtm/forum-api/models/user"
)

var UserPool = sync.Pool{
	New: func() interface{} {
		return &user.User{}
	},
}

var UserArrPool = sync.Pool{
	New: func() interface{} {
		return make(user.Arr, 0, 50)
	},
}
