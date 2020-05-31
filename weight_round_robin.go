package load_balance

import (
	"errors"
	"strconv"
)

type WeightRoundRobinBalance struct {
	curIndex int
	res []*WeightNode
}

type WeightNode struct {
	Addr string
	Weight int
	curWeight int
	effectiveWeight int
}

func (r *WeightRoundRobinBalance) Add(params ...string) error  {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	node := &WeightNode{Addr: params[0], Weight: int(parInt)}
	node.effectiveWeight = node.Weight
	r.res = append(r.res, node)
	return nil
}

func (r *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode
	for i:=0 ; i < len(r.res) ; i++ {
		w := r.res[i]
		// step 1 统计所有的有效权重
		total += w.effectiveWeight
		// step 2 计算当前节点的当前权重值
		w.curWeight += w.effectiveWeight
		// step 3 如果碰到异常情况，需要对有效权重进行 + 1 -1

		// step 4 对比算出有效的权重
		if  best == nil || w.curWeight > best.curWeight {
			best = w
		}
	}

	if best == nil {
		return ""
	}
	// step 5 将当前节点的权重 - 有效权重之和，参与下次循环计算
	best.curWeight -= total
	return best.Addr
}

func (r *WeightRoundRobinBalance) Get(key string) (string, error)  {
	return r.Next(), nil
}