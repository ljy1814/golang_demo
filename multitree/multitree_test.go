package multitree

import (
	"fmt"
	"testing"
)

func TestMTree(t *testing.T) {
	var TestData = []MTreeNode{
		{
			ID:    1,
			PID:   0,
			Value: "梁山泊",
		},
		{
			ID:    2,
			PID:   1,
			Value: "宋江",
		},
		{
			ID:    3,
			PID:   1,
			Value: "卢俊义",
		},
		{
			ID:    4,
			PID:   2,
			Value: "吴用",
		},
		{
			ID:    5,
			PID:   3,
			Value: "公孙胜",
		},
		{
			ID:    6,
			PID:   3,
			Value: "关胜",
		},
		{
			ID:    7,
			PID:   3,
			Value: "林冲",
		},
		{
			ID:    8,
			PID:   2,
			Value: "秦明",
		},
		{
			ID:    9,
			PID:   2,
			Value: "花荣",
		},
		{
			ID:    10,
			PID:   5,
			Value: "蔡进",
		},
	}
	mtp := MergeTree(TestData)
	fmt.Println(mtp)

	for k, v := range mtp {
		fmt.Printf("index[%d], %v\n", k, v)
	}
}
