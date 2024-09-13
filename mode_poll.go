package load_balancer

import (
	"github.com/ameise84/rbmap"
)

func NewBalancerPoll() Balancer {
	mp := rbmap.New[uint64, Context]()
	return &pollBalancer{
		mp:   mp,
		iter: mp.BeginIterator(),
	}
}

type pollBalancer struct {
	mp   rbmap.RBMap[uint64, Context]
	iter rbmap.Iterator[uint64, Context]
}

func (b *pollBalancer) Mode() Mode {
	return PollMode
}

func (b *pollBalancer) Count() int {
	return b.mp.Size()
}

func (b *pollBalancer) Register(ctx Context) error {
	id := ctx.ID()
	if _, ok := b.mp.Load(id); ok {
		return ErrIdRepeat
	}
	b.mp.Store(id, ctx)
	return nil
}

func (b *pollBalancer) UnRegister(id uint64) (v Context, isFind bool) {
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

func (b *pollBalancer) Clean() []Context {
	out := make([]Context, 0, b.mp.Size())
	mp := b.mp.Clean()
	for _, v := range mp {
		out = append(out, v.(Context))
	}
	return out
}

func (b *pollBalancer) Range(f func(Context) bool) {
	b.mp.Range(func(iter rbmap.Iterator[uint64, Context]) bool {
		return f(iter.ValueNoError())
	})
}

func (b *pollBalancer) TakeByID(id uint64) (Context, bool) {
	iter := b.mp.Search(id)
	if iter == b.mp.EndIterator() {
		return nil, false
	}
	return iter.ValueNoError().(Context), true
}

func (b *pollBalancer) TakeNext(uint64) (Context, error) {
	if b.mp.Size() == 0 {
		return nil, ErrIsEmpty
	}
	b.iter, _ = b.iter.Next()
	if b.iter == b.mp.EndIterator() {
		b.iter = b.mp.BeginIterator()
	}
	return b.iter.ValueNoError().(Context), nil
}

func (b *pollBalancer) SetAutoPriWt(int16) {}

func (b *pollBalancer) UpdatePriWt(uint64, int) {}
