package main

import (
	"fmt"
	"math/rand"
	"reflect"
)

type Element interface{}

type LinkNode struct {
	data Element
	next *LinkNode //next指针
}

var listHead *LinkNode

func (this *LinkNode) Insert(nextNode *LinkNode) {
	if listHead == nil {
		listHead = nextNode
		nextNode.next = nil
	}

	var lastNode *LinkNode
	for lastNode = listHead; lastNode.next != nil; lastNode = lastNode.next {

	}
	lastNode.next = nextNode
	nextNode.next = nil
	return
}

func (this *LinkNode) Remove(node *LinkNode) {
	if listHead == nil {
		return
	}
	if listHead.data == node.data {
		listHead = listHead.next
		return
	}

	prevNode := listHead
	if prevNode.next == nil {
		return
	}

	for {
		if prevNode.next.data != node.data {
			prevNode = prevNode.next
			if prevNode == nil || prevNode.next == nil {
				return
			}
		} else {
			break
		}
	}

	prevNode.next = prevNode.next.next
	return
}

func (this *LinkNode) PrintData() {
	for node := listHead; node != nil; node = node.next {
		if reflect.TypeOf(node.data).String() == "int" {
			fmt.Printf("%d ", node.data)
		} else {
			fmt.Println(node.data)
		}
	}
	fmt.Println()
}

func (this *LinkNode) init() {
	for i := 0; i < 20; i++ {
		listHead.Insert(&LinkNode{data: rand.Intn(1024), next: nil})
	}
}

func (this *LinkNode) Traverse() *LinkNode {
	if this == nil || this.next == nil {
		return this
	}
	head := this.next.Traverse()
	this.next.next = this
	this.next = nil
	return head
}

func main() {
	listHead.init()
	listHead.PrintData()
	listHead = listHead.Traverse()
	listHead.PrintData()
}
