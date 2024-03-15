package workflownode

import (
	"fmt"
	"testing"
)

func TestCondition(t *testing.T) {
	var a, b any
	a = true
	b = ""

	b1 := b.(int64)
	fmt.Println(b1)
	if a == b {
		fmt.Printf("%+v, %+v", a, b)
	} else {
		fmt.Printf("%+v, %+v", a, b)
	}
	return
}
