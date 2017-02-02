package search

import "fmt"

type Node struct {
	element interface{} //数据类型
	parent  *Node       //父节点
	left    *Node
	right   *Node
}

//比较器
type comparer func(interface{}, interface{}) bool

//二叉排序树
type Bst struct {
	compare comparer
	root    *Node
}

func (b *Bst) New(compare comparer) *Bst {
	return &Bst{compare: compare}
}

func (b *Bst) inorder_tree_walk(tree *Node) {
	if tree == nil {
		return
	}
	inorder_tree_walk(tree.left)
	fmt.Println(tree.element)
	inorder_tree_walk(tree.right)
}

func (b *Bst) tree_search(tree *Node, element interface{}) *Node {
	if tree == nil {
		return nil
	}
	if tree.element == element {
		return tree
	}
	//左小右大
	if b.compare(element, tree.element) {
		return b.tree_search(tree.right, element)
	} else {
		return b.tree_search(tree.left, element)
	}
}

func (b *Bst) iterative_tree_search(tree *Node, element interface{}) *Node {
	if tree == nil {
		return nil
	}
	for element != tree.element {
		if b.compare(element, tree.element) {
			tree = tree.right
		} else {
			tree = tree.left
		}
	}
	return tree
}

//获取最小元素
func (b *Bst) tree_minimum(tree *Node) *Node {
	for tree != nil {
		tree = tree.left
	}
	return tree
}

//获取最大元素
func (b *Bst) tree_maximum(tree *Node) *Node {
	for tree != nil {
		tree = tree.right
	}
	return tree
}

//后继
func (b *Bst) tree_successor(tree *Node) *Node {
	//有右子树
	if tree.right != nil {
		return b.tree_minimum(tree.right)
	}
	//回溯
	y := tree.parent
	for y != nil && tree == y.right {
		tree = y
		y = y.parent
	}
	return y
}

//前驱
func (b *Bst) tree_predecessor(tree *Node) *Node {
	if tree.left != nil {
		return b.tree_maximum(tree.left)
	}
	y := tree.parent
	for y != nil && tree == y.right {
		tree = y
		y = y.parent
	}
	return y
}

func (b *Bst) tree_insert(tree *Node, element interface{}) *Node {
	y := &Node{}
	for tree != nil {
		//y为双亲节点
		y = tree
		if b.compare(element, tree.element) {
			tree = tree.left
		} else {
			tree = tree.right
		}
	}
	var node *Node = &Node{element, nil, nil}
	node.parent = y

	//空二叉树
	if y == nil {
		b.root = node
		return node
	}

	//找出左边还是右边
	if b.compare(element, y.element) {
		y.left = node
	} else {
		y.right = node
	}
	return y
}

//TODO 以后再仔细处理,此处还要涉及parent指针的设置
func (b *Bst) insert_recursive(tree *Node, element interface{}) *Node {
	if tree == nil {
		var node *Node = &Node{element: element, nil, nil}
		return node
	}
	if b.compare(element, tree.element) {
		cur := insert_recursive(tree.right, element)
		tree.right = cur
		cur.parent = tree
		return cur
	} else {
		cur := insert_recursive(tree.left, element)
		tree.left = cur
		cur.parent = tree
		return cur
	}
}
