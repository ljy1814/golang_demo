package arrayreverse

import (
	"fmt"
	"testing"
)

func TestGetMin(t *testing.T) {
	var arr = []int{3, 4, 5, 6, 7, 1, 2}
	data, index := GetMin(arr)
	fmt.Printf("result : %d, index : %d\n", data, index)

	arr = []int{1, 1, 1, 0, 1, 1}
	data, index = GetMin(arr)
	fmt.Printf("result : %d, index : %d\n", data, index)

	arr = []int{1, 1, 1, 1, 1, 0, 1, 1}
	data, index = GetMin(arr)
	fmt.Printf("result : %d, index : %d\n", data, index)

	arr = []int{1, 1, 2, 3, 0, 1, 1}
	data, index = GetMin(arr)
	fmt.Printf("result : %d, index : %d\n", data, index)

	arr = []int{2, 3, 0, 1, 1}
	data, index = GetMin(arr)
	fmt.Printf("result : %d, index : %d\n", data, index)

	arr = []int{2, 3, 4, 0, 1, 1}
	data, index = GetMin(arr)
	fmt.Printf("result : %d, index : %d\n", data, index)

	arr = []int{2, 3, 4, 0, 1, 1, 1, 1, 1}
	data, index = GetMin(arr)
	fmt.Printf("result : %d, index : %d\n", data, index)
}
