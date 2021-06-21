package util

import "testing"


/*
 IsArraySort Test
 */
func Test_IsArraySort(t *testing.T) {
	type args struct {
		arr   []int
		index int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "正确顺序数组",
			args: args{
				arr:   []int{10, 20, 30, 40, 50, 70},
				index: 6, // 数组长度值
			},
			want: true,
		},
		{
			name: "不正确顺序数组",
			args: args{
				arr:   []int{10, 30, 40, 20, 15, 7},
				index: 6,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsArraySort(tt.args.arr, tt.args.index); got != tt.want {
				t.Errorf("isArraySort() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
  IsArraySort benchmark test
 */
func BenchmarkIsArraySort(b *testing.B) {
	arr := []int{10, 30, 40, 20, 15, 7}
	length := len(arr)
	for i := 0; i < b.N; i++ {
		IsArraySort(arr, length)
	}
}

