package load_balance

import (
	"testing"
)

func TestWeightRoundRobinBalance_Next(t *testing.T) {
	rb := &WeightRoundRobinBalance{}

	rb.Add("127.0.0.1:2004", "4")
	rb.Add("127.0.0.1:2004", "4")
	rb.Add("127.0.0.1:2004", "4")
	rb.Add("127.0.0.1:2004", "4")

	// 定义一个指针及结构体数组，存取所有的地址结构

	// 循环遍历谁的权当前重最大（cur weight = cur weight + effect weight）
	// 权重只是用来作为有限权重的最大值计算

	// 最大的权重，必须减去 所有的有效权重之和，参与到下次计算。当前权重是个动态值，会不断跟进权重进行修复


	// fmt.Printf("res : %v", count)
}
