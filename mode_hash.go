package load_balancer

import (
	"github.com/ameise84/rbmap"
)

func NewBalancerHash() Balancer {
	mp := rbmap.New[uint64, Context]()
	return &hashBalancer{
		mp:   mp,
		iter: mp.BeginIterator(),
	}
}

type hashBalancer struct {
	mp   rbmap.RBMap[uint64, Context]
	iter rbmap.Iterator[uint64, Context]
}

func (b *hashBalancer) Mode() Mode {
	return HashMode
}

func (b *hashBalancer) Count() int {
	return b.mp.Size()
}

func (b *hashBalancer) Register(ctx Context) error {
	id := ctx.ID()
	if _, ok := b.mp.Load(id); ok {
		return ErrIdRepeat
	}
	b.mp.Store(id, ctx)
	return nil
}

func (b *hashBalancer) UnRegister(id uint64) (v Context, isFind bool) {
	key, err := b.iter.Key()
	if err == nil && key == id {
		val := b.iter.ValueNoError()
		b.iter, _ = b.iter.Delete()
		v, isFind = val.(Context), true
	} else {
		if val, ok := b.mp.Delete(id); ok {
			v, isFind = val.(Context), true
		}
	}
	return
}

func (b *hashBalancer) Clean() []Context {
	out := make([]Context, 0, b.mp.Size())
	mp := b.mp.Clean()
	for _, v := range mp {
		out = append(out, v.(Context))
	}
	return out
}

func (b *hashBalancer) Range(f func(Context) bool) {
	b.mp.Range(func(iter rbmap.Iterator[uint64, Context]) bool {
		return f(iter.ValueNoError())
	})
}

func (b *hashBalancer) TakeByID(id uint64) (Context, bool) {
	iter := b.mp.Search(id)
	if iter == b.mp.EndIterator() {
		return nil, false
	}
	return iter.ValueNoError().(Context), true
}

func (b *hashBalancer) TakeNext(id uint64) (Context, error) {
	iter := b.mp.Search(id, rbmap.SearchModeGT|rbmap.SearchModeET)
	if iter == b.mp.EndIterator() {
		if b.mp.Size() > 0 {
			return b.mp.BeginIterator().ValueNoError().(Context), nil
		}
		return nil, ErrIdNotFind
	}
	return iter.ValueNoError().(Context), nil
}

func (b *hashBalancer) SetAutoPriWt(int16) {}

func (b *hashBalancer) UpdatePriWt(uint64, int) {}
