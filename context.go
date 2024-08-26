package load_balancer

import (
	"github.com/ameise84/heap"
	"github.com/ameise84/heap/compare"
)

type contextWrap struct {
	ctx    Context
	id     heap.ID
	weight int //权重值越大,优先级越低.最大值为maxInt,最小值为0
}

func (i *contextWrap) Compare(other compare.Ordered) compare.Result {
	d := i.weight - other.(*contextWrap).weight
	if d < 0 {
		return compare.Smaller
	}
	if d > 0 {
		return compare.Larger
	}
	return compare.Equal
}
