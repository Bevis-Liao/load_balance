package load_balance

import (
	"errors"
	"math/rand"
)

type RandomBalance struct {
	curIdx int
	res []string
	// 观察主体
}

func (rb *RandomBalance) Add(params ...string) error {
	length := len(params)
	if length == 0 {
		return errors.New("param len 1 at least")
	}
	for i:=0; i < length; i++{
		rb.res = append(rb.res, params[i])
	}
	return nil
}

func (rb *RandomBalance) Next() string {
	length := len(rb.res)
	if length == 0 {
		return ""
	}

	idx := rand.Intn(length)
	return rb.res[idx]
}

func (rb *RandomBalance) Get(key string) (string, error) {
	return rb.Next(), nil
}