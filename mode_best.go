package load_balancer

import (
	"github.com/ameise84/heap"
	"log"
	"math"
)

func NewBalancerBest() Balancer {
	return &bestBalancer{
		h:             heap.NewIDHeapMin[*contextWrap](64),
		mp:            make(map[uint64]struct{}),
		autoAddWeight: 1,
	}
}

type bestBalancer struct {
	h             heap.IDHeap[*contextWrap]
	mp            map[uint64]struct{}
	autoAddWeight int //当获取时,自动添加权重值,默认1
}

func (b *bestBalancer) Mode() Mode {
	return BestMode
}

func (b *bestBalancer) Count() int {
	return len(b.mp)
}

func (b *bestBalancer) Register(ctx Context) error {
	id := ctx.ID()
	i := &contextWrap{ctx: ctx, id: id, weight: 0} //注册时默认最优
	err := b.h.Push(id, i)
	if err != nil {
		return err
	}
	b.mp[id] = struct{}{}
	return nil
}

func (b *bestBalancer) UnRegister(id uint64) (Context, bool) {
	if _, ok := b.mp[id]; ok {
		v, _ := b.h.Remove(id)
		if v.weight > 0 {
			log.Printf("bestBalancer UnRegister -> Id[%d] is using\n", id)
		}
		delete(b.mp, id)
		return v.ctx, true
	}
	return nil, false
}

func (b *bestBalancer) Clean() []Context {
	out := make([]Context, 0, len(b.mp))
	vs := b.h.CleanToSlice()
	for _, v := range vs {
		out = append(out, v.ctx)
	}
	b.mp = map[uint64]struct{}{}
	return out
}

func (b *bestBalancer) Range(f func(Context) bool) {
	b.h.Range(func(v *contextWrap) bool {
		return f(v.ctx)
	})
}

func (b *bestBalancer) TakeByID(id uint64) (Context, bool) {
	if v, ok := b.h.Find(id); ok {
		v.weight += b.autoAddWeight
		if v.weight < 0 {
			v.weight = math.MaxInt
		}
		b.h.Update(v.id, v)
		return v.ctx, true
	}
	return nil, false
}

func (b *bestBalancer) TakeNext(uint64) (Context, error) {
	_, v, err := b.h.Peek()
	if err != nil {
		return nil, err
	}
	v.weight += b.autoAddWeight
	b.h.Update(v.id, v)
	return v.ctx, nil
}

func (b *bestBalancer) SetAutoPriWt(weight uint16) {
	b.autoAddWeight = int(weight)
}

func (b *bestBalancer) UpdatePriWt(id uint64, weight int) {
	if v, ok := b.h.Find(id); ok {
		if weight < 0 {
			weight = 0
		}
		v.weight = weight
		b.h.Update(id, v)
	}
}
