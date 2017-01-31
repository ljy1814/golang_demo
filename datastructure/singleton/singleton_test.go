package singleton

import (
	"fmt"
	"testing"
)

func TestS1(t *testing.T) {
	fmt.Println(GetInstance())
	fmt.Printf("addr : %p\n", GetInstance())
	fmt.Println(GetInstance())
	fmt.Printf("addr : %p\n", GetInstance())
	fmt.Println(GetInstance())
	fmt.Printf("addr : %p\n", GetInstance())
}
func TestS2(t *testing.T) {
	fmt.Println(GetInstance2())
	fmt.Printf("addr : %p\n", GetInstance2())
	fmt.Println(GetInstance2())
	fmt.Printf("addr : %p\n", GetInstance2())
	fmt.Println(GetInstance2())
	fmt.Printf("addr : %p\n", GetInstance2())
}
func TestS3(t *testing.T) {
	fmt.Println(GetInstance3())
	fmt.Printf("addr : %p\n", GetInstance3())
	fmt.Println(GetInstance3())
	fmt.Printf("addr : %p\n", GetInstance3())
	fmt.Println(GetInstance3())
	fmt.Printf("addr : %p\n", GetInstance3())
}
