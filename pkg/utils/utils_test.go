package utils

import (
	"fmt"
	"testing"
)

type TestType int

const (
	Val1 TestType = iota
	Val2
	Val3
)

func TestFillMap(t *testing.T) {
	m := make(map[int32]TestType)
	var rows, cols int32 = 5, 3
	FillMap(m, rows*cols, Val2)
	fmt.Printf("size:%dx%d\n", rows, cols)
	//keys := reflect.ValueOf(&m).Elem().MapKeys()
	for k := range rows * cols {
		fmt.Printf("key:%v -> (x,y)=(%d, %d)\n", k, k/cols, k%cols)
	}

}
