package example

import (
	"fmt"
	"github.com/ameise84/load_balancer"
	"math/rand/v2"
	"sort"
	"testing"
)

type item struct {
	id uint64
}

func (i item) ID() uint64 {
	return i.id
}

func TestBalancer(t *testing.T) {
	testPoll()
	fmt.Println("=====================")
	testBest()
	fmt.Println("=====================")
	testHash()
}

func testHash() {
	lb := load_balancer.NewBalancerHash()
	ls := make([]*item, 0, 64)
	for i := uint64(0); i < 64; i++ {
		ls = append(ls, &item{rand.Uint64N(256)})
	}
	for i := uint64(0); i < 64; i++ {
		_ = lb.Register(ls[i])
	}
	finds := make([]uint64, 0, 64)
	for i := uint64(0); i < 64; i++ {
		finds = append(finds, rand.Uint64N(256))
	}
	sort.Slice(finds, func(i, j int) bool { return finds[i] < finds[j] })
	for i := uint64(0); i < 64; i++ {
		find := finds[i]
		x, err := lb.TakeNext(find)
		if err != nil {
			return
		}
		fmt.Println(find, " --> ", x.ID())
	}
}

func testBest() {
	lb := load_balancer.NewBalancerBest()
	ls := make([]*item, 0, 64)
	for i := uint64(0); i < 64; i++ {
		ls = append(ls, &item{i})
	}
	rand.Shuffle(len(ls), func(i, j int) {
		ls[i], ls[j] = ls[j], ls[i]
	})
	for i := uint64(0); i < 32; i++ {
		_ = lb.Register(ls[i])
	}
	for i := uint64(0); i < 10; i++ {
		x, err := lb.TakeNext(0)
		if err != nil {
			return
		}
		fmt.Println(x.ID())
	}
	fmt.Println("-------------")
	for i := uint64(32); i < 64; i++ {
		_ = lb.Register(ls[i])
	}
	for i := uint64(0); i < 64; i++ {
		x, err := lb.TakeNext(0)
		if err != nil {
			return
		}
		fmt.Println(x.ID())
	}
}

func testPoll() {
	lb := load_balancer.NewBalancerPoll()
	ls := make([]*item, 0, 64)
	for i := uint64(0); i < 64; i++ {
		ls = append(ls, &item{i})
	}
	rand.Shuffle(len(ls), func(i, j int) {
		ls[i], ls[j] = ls[j], ls[i]
	})
	for i := uint64(0); i < 32; i++ {
		_ = lb.Register(ls[i])
	}

	for i := uint64(0); i < 10; i++ {
		x, err := lb.TakeNext(0)
		if err != nil {
			return
		}
		fmt.Println(x.ID())
	}
	fmt.Println("===")
	for i := uint64(32); i < 64; i++ {
		_ = lb.Register(ls[i])
	}
	for i := uint64(0); i < 64; i++ {
		x, err := lb.TakeNext(0)
		if err != nil {
			return
		}
		fmt.Println(x.ID())
	}
}
