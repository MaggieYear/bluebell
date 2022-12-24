package models

// 结构体内存堆积概念
import (
	"fmt"
	"testing"
	"unsafe"
)

type s1 struct {
	a int8
	b string
	c int8
}

type s2 struct {
	a int8
	b int8
	c string
}

func TestStruct(t *testing.T) {
	v1 := s1{a: 1, b: "x", c: 2}
	v2 := s2{a: 1, b: 2, c: "x"}

	// 对比两个结构体的大小
	fmt.Println(unsafe.Sizeof(v1)) // 32
	fmt.Println(unsafe.Sizeof(v2)) // 24
}
