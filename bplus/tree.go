package bplus

import "io"

type Tree struct {
	c     int
	cmp   Cmp
	first *d
	last  *d
	r     interface{}
	ver   int64
}

// return a new empty tree
func NewTree(cmp Cmp) *Tree {
	return btTPool.get(cmp)
}

// removes all k/v pairs from tree
func (t *Tree) Clear() {
	if t.r == nil {
		return
	}
	clr(t.r)
	t.c, t.first, t.last, t.r = 0, nil, nil, nil
	t.ver++
}

// performs clear; recycles t (pool) for later reuse if needed.
// no refrences to t should exist.
func (t *Tree) Close() {
	t.Clear()
	*t = zt
}

// remove k/v pair if exists, return true; else false
func (t *Tree) Delete(k []byte) (ok bool) {
	pi := -1
	var p *x
	q := t.r
	if q == nil {
		return false
	}
	for {
		var i int
		i, ok = t.find(q, k)
		if ok {
			switch x := q.(type) {
			case *x:
				if x.c < kx && q != t.r {
					x, i = t.underflowX(p, x, pi, i)
				}
				pi = i + 1
				p = x
				q = x.x[pi].ch
				ok = false
				continue
			case *d:
				t.extract(x, i)
				if x.c >= kd {
					return true
				}

				if q != t.r {
					t.underflow(p, x, pi)
				} else if t.c == 0 {
					t.Clear()
				}
				return true
			}
		}
		switch x := q.(type) {
		case *x:
			if x.c < kx && q != t.r {
				x, i = t.underflowX(p, x, pi, i)
			}
			pi = i
			p = x
			q = x.x[i].ch
		case *d:
			return false
		}
	}
}

// return first item in tree in lexigraphical order...
// or zero value if tree is empty
func (t *Tree) First() (k []byte, v []byte) {
	if q := t.first; q != nil {
		q := &q.d[0]
		k, v = q.k, q.v
	}
	return
}

// return value of k and true if exists; else return zero value and false
func (t *Tree) Get(k []byte) (v []byte, ok bool) {
	q := t.r
	if q == nil {
		return
	}
	for {
		var i int
		if i, ok = t.find(q, k); ok {
			switch x := q.(type) {
			case *x:
				q = x.x[i+1].ch
				continue
			case *d:
				return x.d[i].v, true
			}
		}
		switch x := q.(type) {
		case *x:
			q = x.x[i].ch
		default:
			return
		}
	}
}

func (t *Tree) Last() (k []byte, v []byte) {
	if q := t.last; q != nil {
		q := &q.d[q.c-1]
		k, v = q.k, q.v
	}
	return
}

func (t *Tree) Len() int {
	return t.c
}

func (t *Tree) Put(k []byte, upd func(oldV []byte, exists bool) (newV []byte, write bool)) (oldV []byte, written bool) {
	pi := -1
	var p *x
	q := t.r
	var newV []byte
	if q == nil {
		// new KV pair in empty tree
		newV, written = upd(newV, false)
		if !written {
			return
		}

		z := t.insert(btDPool.Get().(*d), 0, k, newV)
		t.r, t.first, t.last = z, z, z
		return
	}
	for {
		i, ok := t.find(q, k)
		if ok {
			switch x := q.(type) {
			case *x:
				if x.c > 2*kx {
					x, i = t.splitX(p, x, pi, i)
				}
				pi = i + 1
				p = x
				q = x.x[i+1].ch
				continue
			case *d:
				oldV = x.d[i].v
				newV, written = upd(oldV, true)
				if !written {
					return
				}
				x.d[i].v = newV
			}
			return
		}
		switch x := q.(type) {
		case *x:
			if x.c > 2*kx {
				x, i = t.splitX(p, x, pi, i)
			}
			pi = i
			p = x
			q = x.x[i].ch
		case *d: // new KV pair
			newV, written = upd(newV, false)
			if !written {
				return
			}
			switch {
			case x.c < 2*kd:
				t.insert(x, i, k, newV)
			default:
				t.overflow(p, x, pi, i, k, newV)
			}
			return
		}
	}
}

// return and Enumerator postion on an item where k >= item's key
func (t *Tree) Seek(k []byte) (e *Enumerator, ok bool) {
	q := t.r
	if q == nil {
		e = btEPool.get(nil, false, 0, k, nil, t, t.ver)
		return
	}
	for {
		var i int
		if i, ok = t.find(q, k); ok {
			switch x := q.(type) {
			case *x:
				q = x.x[i+1].ch
				continue
			case *d:
				return btEPool.get(nil, ok, i, k, x, t, t.ver), true
			}
		}
		switch x := q.(type) {
		case *x:
			q = x.x[i].ch
		case *d:
			return btEPool.get(nil, ok, i, k, x, t, t.ver), false
		}
	}
}

// returns an enumerator positioned on the first k/v pair in the tree
func (t *Tree) SeekFirst() (e *Enumerator, err error) {
	q := t.first
	if q == nil {
		return nil, io.EOF
	}
	return btEPool.get(nil, true, 0, q.d[0].k, q, t, t.ver), nil
}

// returns an enumerator positioned on the last k/v pair in the tree
func (t *Tree) SeekLast() (e *Enumerator, err error) {
	q := t.last
	if q == nil {
		return nil, io.EOF
	}
	return btEPool.get(nil, true, q.c-1, q.d[q.c-1].k, q, t, t.ver), nil
}

// sets the value associated with k
func (t *Tree) Set(k []byte, v []byte) {
	//dbg("--- PRE Set(%v, %v)\n%s", k, v, t.dump())
	//defer func() {
	//  dbg("--- POST\n%s\n====\n", t.dump())
	//}()
	pi := -1
	var p *x
	q := t.r
	if q == nil {
		z := t.insert(btDPool.Get().(*d), 0, k, v)
		t.r, t.first, t.last = z, z, z
		return
	}
	for {
		i, ok := t.find(q, k)
		if ok {
			switch x := q.(type) {
			case *x:
				if x.c > 2*kx {
					x, i = t.splitX(p, x, pi, i)
				}
				pi = i + 1
				p = x
				q = x.x[i+1].ch
				continue
			case *d:
				x.d[i].v = v
			}
			return
		}
		switch x := q.(type) {
		case *x:
			if x.c > 2*kx {
				x, i = t.splitX(p, x, pi, i)
			}
			pi = i
			p = x
			q = x.x[i].ch
		case *d:
			switch {
			case x.c < 2*kd:
				t.insert(x, i, k, v)
			default:
				t.overflow(p, x, pi, i, k, v)
			}
			return
		}
	}
}
