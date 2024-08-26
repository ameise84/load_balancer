package load_balancer

import "errors"

var (
	ErrIdRepeat  = errors.New("balancer join id repeat")
	ErrIdNotFind = errors.New("load_balancer take not find id")
	ErrIsEmpty   = errors.New("load_balancer is empty")
)
