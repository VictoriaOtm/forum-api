package userstore

import "sync"

type userPool struct {
	sync.Pool
}

func (p *userPool) Acquire() *User {
	return p.Get().(*User)
}

func (p *userPool) Utilize(f *User) {
	p.Put(f)
}

var Pool = userPool{
	sync.Pool{
		New: func() interface{} {
			return &User{}
		},
	},
}

type userSlicePool struct {
	sync.Pool
}

func (p *userSlicePool) Acquire() UserSlice {
	return p.Get().(UserSlice)
}

func (p *userSlicePool) Utilize(f UserSlice) {
	p.Put(f[:0])
}

var PoolUserSlice = userSlicePool{
	sync.Pool{
		New: func() interface{} {
			return make(UserSlice, 0, 100)
		},
	},
}
