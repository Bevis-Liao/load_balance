package load_balance

import (
	"errors"
)

type RoundRobinBalance struct {
	curIdx int
	res []string
}

func (rb *RoundRobinBalance) Add(params ...string) error {
	length := len(params)
	if length == 0 {
		return errors.New("param len 1 at least")
	}
	for i:=0; i < length; i++{
		rb.res = append(rb.res, params[i])
	}
	return nil
}

func (rb *RoundRobinBalance) Next() string {
	length := len(rb.res)
	if length == 0 {
		return ""
	}

	if rb.curIdx >= (length - 1) {
		rb.curIdx = 0
	} else {
		rb.curIdx++
	}

	return rb.res[rb.curIdx]
}

func (rb *RoundRobinBalance) Get(key string) (string, error)  {
	return rb.Next(), nil
}