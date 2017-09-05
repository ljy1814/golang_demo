package btree

import (
	"os"
	"testing"
)

type item int

func (i item) Less(other Item) bool {
	return i < other.(item)
}

func TestGet(t *testing.T) {
	bt := New(2)

	bt.Insert(item(4))
	bt.Print(os.Stdout)
	bt.Insert(item(8))
	bt.Print(os.Stdout)
	bt.Insert(item(18))
	bt.Print(os.Stdout)
	bt.Insert(item(20))
	bt.Print(os.Stdout)
	bt.Insert(item(22))
	bt.Print(os.Stdout)
	bt.Insert(item(24))
	bt.Print(os.Stdout)
	bt.Insert(item(26))
	bt.Print(os.Stdout)
	bt.Insert(item(28))
	bt.Print(os.Stdout)
	bt.Insert(item(30))
	bt.Print(os.Stdout)

	if bt.Get(item(11)) != nil {
		t.Fatalf("expected is nil, actual is %v", bt.Get(item(11)))
	}
}
