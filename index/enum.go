package index

import "io"

type Enumerator struct {
	err error
	hit bool
	i   int
	k   int64
	q   *d
	t   *Tree
	ver int64
}

// recycles e (pool) for possible reuse.
// no references to e should exists.
func (e *Enumerator) Close() {
	*e = ze
	btEPool.Put(e)
}

// returns the currently enumerated item, if it exists...
// and moves to the next item in the key order.
func (e *Enumerator) Next() (k int64, v []byte, err error) {
	if err = e.err; err != nil {
		return
	}
	if e.ver != e.t.ver {
		f, hit := e.t.Seek(e.k)
		if !e.hit && hit {
			if err = f.next(); err != nil {
				return
			}
		}
		*e = *f
		f.Close()
	}
	if e.q == nil {
		e.err, err = io.EOF, io.EOF
		return
	}
	if e.i >= e.q.c {
		if err = e.next(); err != nil {
			return
		}
	}
	i := e.q.d[e.i]
	k, v = i.k, i.v
	e.k, e.hit = k, false
	e.next()
	return
}

// returns the currently enumerated item, it it exists...
// and moves to the previous item in the key order.
func (e *Enumerator) Prev() (k int64, v []byte, err error) {
	if err = e.err; err != nil {
		return
	}
	if e.ver != e.t.ver {
		f, hit := e.t.Seek(e.k)
		if !e.hit && hit {
			if err = f.prev(); err != nil {
				return
			}
		}
		*e = *f
		f.Close()
	}
	if e.q == nil {
		e.err, err = io.EOF, io.EOF
		return
	}
	if e.i >= e.q.c {
		if err = e.next(); err != nil {
			return
		}
	}
	i := e.q.d[e.i]
	k, v = i.k, i.v
	e.k, e.hit = k, false
	e.prev()
	return
}
