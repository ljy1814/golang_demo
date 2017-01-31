package main

import "fmt"

//右边的数都比左边的大,下面的数都比上面的大

var (
	data [4][4]int
	len1 = 4
	len2 = 4
)

func init() {
	data = [4][4]int{
		{1, 2, 8, 9},
		{2, 4, 9, 12},
		{4, 7, 10, 13},
		{6, 8, 11, 15},
	}
}
func find(key int) (i, j int) {
	for i := 0; i < len1; {
		for j := len2 - 1; j >= 0; {
			if i >= len1 || j < 0 {
				break
			}
			if data[i][j] == key {
				return i, j
			} else if data[i][j] > key {
				j -= 1
			} else {
				i++
			}
		}
	}
	return -1, -1
}
func main() {
	fmt.Println("vim-go")
	i, j := find(10)
	fmt.Printf("%d %d\n", i, j)
	i, j = find(15)
	fmt.Printf("%d %d\n", i, j)
	i, j = find(17)
	fmt.Printf("%d %d\n", i, j)
	i, j = find(14)
	fmt.Printf("%d %d\n", i, j)
	i, j = find(7)
	fmt.Printf("%d %d\n", i, j)
}
