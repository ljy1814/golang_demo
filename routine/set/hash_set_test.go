package set

import (
	"fmt"
	"runtime/debug"
	"errors"
//	"strings"
	"testing"
)

func TestHashSetCreateion(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	t.Log("Starting TestHashSetCreateion...")
	hs := NewHashSet()
	t.Logf("Create a HashSet value : %v\n", hs)
	if hs == nil {
		t.Errorf("The result of func NewHashSet is nil!\n")
	}
	isSet := IsSet(hs)
	if !isSet {
		t.Errorf("The value of HashSet is not Set!\n")
	} else {
		t.Logf("The HashSet value is a Set.\n")
	}
}
/*
func TestHashSetLenAndContains(t *testing.T) {
	testSetLenAndContains(t, func() Set { return NewHashSet() }, "HashSet")
}

func TestHashSetAdd(t *testing.T) {
	testSetAdd(t, func() Set { return NewHashSet() }, "HashSet")
}

func TestHashSetRemove(t *testing.T) {
	testSetRemove(t, func() Set { return NewHashSet() }, "HashSet")
}

func TestHashSetClear(t *testing.T) {
	testSetClear(t, func() Set { return NewHashSet() }, "HashSet")
}

func TestHashSetElements(t *testing.T) {
	testSetElements(t, func() Set { return NewHashSet() }, "HashSet")
}

func TestHashSetSame(t *testing.T) {
	testSetSame(t, func() Set { return NewHashSet() }, "HashSet")
}

func TestSetString(t *testing.T) {
	testSetString(t, func() Set { return NewHashSet() }, "HashSet")
}

func testSetOp(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %d\n", err)
		}
	}()
	fmt.Println(222)
	t.Logf("Starting TestHashSetOp...")
	hs := NewHashSet()
	if hs.Len() != 0 {
		t.Errorf("ERROR: The length of original HashSet value is not 0!\n")
		t.FailNow()
	}
	randElem := genRandElement()
	expecetedElemMap := make(map[interface{}]bool)
	t.Logf("Add %v to the HashSet value %v.\n", randElem, hs)
	hs.Add(randElem)
	expectedElemMap[randElem] = true
	expectedLen := len(exceptedElemMap)
	if hs.Len() != exceptedLen {
		t.Errorf("ERROR: The length of HashSet value %d is not %d!.\n", hs.Len(), expectedLen)
	}
	var result bool
	for i := 0; i < 8; i++ {
		randElem = genRandElem()
		t.Logf("Add %v to the HashSet value %v.\n"), randElem, hs)
		result = hs.Add(randElem)
	}
}
*/

func TestErrors(t *testing.T) {
	fmt.Println("errors testing....")
	err := errors.New("haha")
	fmt.Println(err)
	fmt.Println(errors.New("haha"))
}
