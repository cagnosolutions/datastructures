package bplus

import "sync"

var (
	btDPool = sync.Pool{
		New: func() interface{} {
			return &d{}
		},
	}

	btEPool = btEpool{
		sync.Pool{
			New: func() interface{} {
				return &Enumerator{}
			},
		},
	}

	btTPool = btTpool{
		sync.Pool{
			New: func() interface{} {
				return &Tree{}
			},
		},
	}

	btXPool = sync.Pool{
		New: func() interface{} {
			return &x{}
		},
	}
)

type btTpool struct {
	sync.Pool
}

func (p *btTpool) get(cmp Cmp) *Tree {
	x := p.Get().(*Tree)
	x.cmp = cmp
	return x
}

type btEpool struct {
	sync.Pool
}

func (p *btEpool) get(err error, hit bool, i int, k []byte, q *d, t *Tree, ver int64) *Enumerator {
	x := p.Get().(*Enumerator)
	x.err, x.hit, x.i, x.k, x.q, x.t, x.ver = err, hit, i, k, q, t, ver
	return x
}
