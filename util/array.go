package util

import (
	"fmt"
	"reflect"
)

/*
数组元素是否有序，升序
@param arr 待检查数组
@param length 数组长度
@return true 为是升序
*/
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

/*
数组 array 是否包含 val【reflect实现】
@param array 数组
@param val 数组元素
@return index 数组元素的索引
*/
func Contains(array interface{}, val interface{}) (index int) {
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		{
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(val, s.Index(i).Interface()) {
					index = i
					return
				}
			}
		}
	}
	return
}

/*
字符串数组 array 是否包含 val【for实现】
@param array 数组
@param val 数组元素
@return index 数组元素的索引
*/
func StringsContains(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

/*
整形数组 array 是否包含 val【for实现】
@param array 数组
@param val 数组元素
@return index 数组元素的索引
*/
func IntContains(array []int, val int) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}
