package model_pool

import (
	"sync"
	"github.com/VictoriaOtm/forum-api/models/thread"
)

var ThreadPool = sync.Pool{
	New: func() interface{} {
		return &thread.Thread{}
	},
}

var ThreadArrPool = sync.Pool{
	New: func() interface{} {
		return make(thread.Arr, 0, 300)
	},
}
