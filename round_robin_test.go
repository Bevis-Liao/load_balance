package load_balance

import (
	"fmt"
	"testing"
)

func TestRoundRobinBalance_Next(t *testing.T) {
	rb := &RoundRobinBalance{}
	rb.Add("127.0.0.1:2003",
		"127.0.0.1:2004",
		"127.0.0.1:2005",
		"127.0.0.1:2006",
	)

	for i:=0; i < 20; i++ {
		fmt.Printf("Get access ip : %s \n", rb.Next())
	}
}
