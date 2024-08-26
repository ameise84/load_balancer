package load_balancer

type Mode int

const (
	RandMode Mode = iota //随机
	PollMode             //轮询
	BestMode             //最优
	HashMode             //一致性hash
)

type Context interface {
	ID() uint64
}

type Balancer interface {
	Mode() Mode
	Count() int // register count
	Register(ctx Context) error
	UnRegister(ctxID uint64) (Context, bool)
	Clean() []Context
	Range(func(Context) bool)              // Scan Context when return true
	TakeByID(ctxID uint64) (Context, bool) // return context with ID
	TakeNext(id uint64) (Context, error)   // RandMode , BestMode and PollMode ignore id , HashMode use id on hash ring
	SetAutoPriWt(weight uint16)            //BestMode set auto add  priority weight , each time take. default value is 1
	UpdatePriWt(ctxID uint64, weight int)  // BestMode the priority weight change, each time take, the priority weight value add by [AutoPriWt] auto. the lower the weight value, the higher the priority. the valid range is [0,maxInt]
}

func NewBalancer(m Mode) Balancer {
	switch m {
	case PollMode:
		return NewBalancerPoll()
	case BestMode:
		return NewBalancerBest()
	case HashMode:
		return NewBalancerHash()
	default:
		panic("new balancer ill mode")
	}
	return nil
}
