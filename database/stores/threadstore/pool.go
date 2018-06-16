package threadstore

import "sync"

type threadPool struct {
	sync.Pool
}

func (p *threadPool) Acquire() *Thread {
	return p.Get().(*Thread)
}

func (p *threadPool) Utilize(t *Thread) {
	t.Slug = nil
	t.Votes = 0
	p.Put(t)
}

var Pool = threadPool{
	sync.Pool{
		New: func() interface{} {
			return &Thread{}
		},
	},
}

type threadSlicePool struct {
	sync.Pool
}

func (p *threadSlicePool) Acquire() ThreadSlice {
	return p.Get().(ThreadSlice)
}

func (p *threadSlicePool) Utilize(ts ThreadSlice) {
	p.Put(ts[:0])
}

var PoolThreadSlice = threadSlicePool{
	sync.Pool{
		New: func() interface{} {
			return make(ThreadSlice, 0, 100)
		},
	},
}
