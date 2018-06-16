package forumstore

import "sync"

type forumPool struct {
	sync.Pool
}

func (p *forumPool) Acquire() *Forum {
	return p.Get().(*Forum)
}

func (p *forumPool) Utilize(f *Forum) {
	f.Posts = 0
	f.Threads = 0
	p.Put(f)
}

var Pool = forumPool{
	sync.Pool{
		New: func() interface{} {
			return &Forum{}
		},
	},
}
