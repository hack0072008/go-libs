package util

import "fmt"

// 数组元素是否有序，升序
// @param arr 待检查数组
// @param length 数组长度
// @return true 为是升序
func IsArraySort(arr []int, length int) bool {
	if length == 1 {
		return true
	} else {
		if arr[length-1] <= arr[length-2] {
			return false
		} else {
			return IsArraySort(arr, length-1) // 升序
		}
	}
}

/*
 example
 */
func ExampleIsPalindrome() {
	arr := []int{10, 30, 40, 20, 15, 7}
	l := len(arr)
	fmt.Println(IsArraySort(arr, l))
	// Output:
	// true
	// false
}

