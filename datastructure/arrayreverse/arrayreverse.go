package arrayreverse

/*
	处理数组如{3,4,5,6,7,8,1,2} //找到1
	{1,1,1,0,1,1,1} 找到0
	循环移位之后的有序数组,最左边大于等于最右边
*/
func GetMin(arr []int) (result, index int) {
	mid := 0
	for left, right := 0, len(arr)-1; arr[left] >= arr[right]; {
		mid = (left + right) >> 1
		if right-left == 1 {
			mid = right
			break
		}
		if arr[mid] == arr[left] && arr[mid] == arr[right] {
			return getMin(arr)
		}
		if arr[mid] >= arr[right] {
			left = mid
		} else if arr[mid] <= arr[right] {
			right = mid
		}
	}
	return arr[mid], mid
}

func getMin(arr []int) (result, index int) {
	result = arr[0]
	index = 0
	for i := 1; i < len(arr); i++ {
		if arr[i] < result {
			result = arr[i]
			index = i
		}
	}
	return
}
