package poststore

import "sync"

type postSlicePool struct {
	sync.Pool
}

func (p *postSlicePool) Acquire() PostSlice {
	return p.Get().(PostSlice)
}

func (p *postSlicePool) Utilize(ts PostSlice) {
	p.Put(ts[:0])
}

var PoolPostSlice = postSlicePool{
	sync.Pool{
		New: func() interface{} {
			return make(PostSlice, 0, 100)
		},
	},
}
