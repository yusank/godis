package datastruct

import (
	"fmt"
	"strings"
)

type List struct {
	length     int
	head, tail *listNode
}

type listNode struct {
	next, prev *listNode
	value      string
}

func newListByLPush(values ...string) *List {
	list := new(List)
	for _, v := range values {
		list.LPush(v)
	}

	return list
}

func newListByRPush(values ...string) *List {
	list := new(List)
	for _, v := range values {
		list.RPush(v)
	}

	return list
}

func newListNode(val string) *listNode {
	return &listNode{
		value: val,
	}
}

func (n *listNode) addToHead(val string) *listNode {
	node := newListNode(val)
	n.prev = node
	node.next = n

	return node
}

// pop current node and return next node
func (n *listNode) popAndNext() *listNode {
	var next = n.next

	n.next = nil
	if next != nil {
		next.prev = nil
	}

	return next
}

// pop current node and return prev node
func (n *listNode) popAndPrev() *listNode {
	var prev = n.prev

	n.prev = nil
	if prev != nil {
		prev.next = nil
	}

	return prev
}

func (n *listNode) addToTail(val string) *listNode {
	node := newListNode(val)
	n.next = node
	node.prev = n

	return node
}

func (l *List) LPush(val string) {
	l.length++
	if l.head != nil {
		l.head = l.head.addToHead(val)
		return
	}

	node := newListNode(val)
	l.head = node
	l.tail = node
}

func (l *List) LPop() (val string, ok bool) {
	if l.head == nil {
		return "", false
	}

	l.length--
	val = l.head.value
	l.head = l.head.popAndNext()
	if l.head == nil {
		l.tail = nil
	}

	return val, true
}

func (l *List) RPush(val string) {
	l.length++
	if l.tail != nil {
		l.tail = l.tail.addToTail(val)
		return
	}

	node := newListNode(val)
	l.head = node
	l.tail = node
}

func (l *List) RPop() (val string, ok bool) {
	if l.tail == nil {
		return "", false
	}

	l.length--
	val = l.tail.value
	l.tail = l.tail.popAndPrev()
	if l.tail == nil {
		l.head = nil
	}

	return val, true
}

func (l *List) LRange(start, stop int) (values []string) {
	if l.head == nil {
		return nil
	}

	if start < 0 {
		start = start + l.length
		if start < 0 {
			start = 0
		}
	}

	if stop < 0 {
		stop = stop + l.length
	}

	// start already >=0 , so if stop < 0 then this case is true
	if start > stop || start > l.length {
		return nil
	}
	if stop >= l.length {
		stop = l.length - 1
	}

	var (
		head = l.head
		idx  int
	)

	for head != nil && idx <= stop {
		if idx >= start {
			values = append(values, head.value)
		}

		idx++
		head = head.next
	}

	return
}

// LRemAll .
// count == 0 => remove all element equal value from head to tail
// count > 0 => remove count element equal value from head to tail
// count < 0 => remove count element equal value from tail to head

// LRemCountFromHead n == l.length if remove all
func (l *List) LRemCountFromHead(value string, n int) (cnt int) {
	var (
		dumbHead = &listNode{
			next: l.head,
		}
		prev = dumbHead
		cur  = l.head
		next = l.head.next
	)

	for cur != nil && n > 0 {
		if cur.value == value {
			prev.next = next
			if next != nil {
				next.prev = prev
			}
			// drop  cur node
			cur.prev, cur.next = nil, nil

			cnt++
			n--
		} else {
			// only  move when cur node not removed
			prev = prev.next
		}

		cur = next
		if next != nil {
			next = next.next
		}
	}

	// remove last element
	if prev.next == nil {
		l.tail = prev
	}

	l.head = dumbHead.next
	if l.head != nil {
		// if remove first element from, dumbHead.next.prev will be point to dumbHead
		// // In other words, l.head.tail will be not nil
		l.head.prev = nil
	}
	l.length -= cnt
	return
}

// LRemCountFromTail n is positive for convenience
func (l *List) LRemCountFromTail(value string, n int) (cnt int) {
	var (
		dumbTail = &listNode{
			prev: l.tail,
		}
		prev = dumbTail
		cur  = l.tail
		next = l.tail.prev
		// 1   2   3    4     5    dumTail
		//              ^     ^      ^  <<-
		//            next   cur    prev
	)

	// range from tail to head
	for cur != nil && n > 0 {
		if cur.value == value {
			prev.prev = next
			if next != nil {
				next.next = prev
			}
			// drop  cur node
			cur.prev, cur.next = nil, nil

			cnt++
			n--
		} else {
			// only  move when cur node not removed
			prev = prev.prev
		}

		cur = next
		if next != nil {
			next = next.prev
		}
	}

	// remove last(head) element
	if prev.prev == nil {
		l.head = prev
	}

	l.tail = dumbTail.prev
	if l.tail != nil {
		// if remove first element from tail, dumbTail.prev.next will be point to dumbTail
		// In other words, l.tail.next will be not nil
		l.tail.next = nil
	}

	l.length -= cnt
	return
}

func (l *List) lIndexNode(i int) *listNode {
	if l.head == nil {
		return nil
	}

	if i < 0 {
		i += l.length
	}

	if i >= l.length || i < 0 {
		return nil
	}

	var (
		idx     int
		cur     = l.head
		reverse = i > l.length/2+1
	)

	if reverse {
		idx = l.length - 1
		cur = l.tail
	}

	for i != idx {
		if reverse {
			idx--
			cur = cur.prev
		} else {
			idx++
			cur = cur.next
		}
	}

	return cur
}

func (l *List) LIndex(i int) (string, bool) {
	node := l.lIndexNode(i)
	if node == nil {
		return "", false
	}

	return node.value, true
}

func (l *List) LSet(i int, val string) bool {
	node := l.lIndexNode(i)
	if node == nil {
		return false
	}

	node.value = val
	return true
}

func (l *List) LInsert(target, newValue string, flag int) bool {
	if l.head == nil {
		return false
	}

	node := l.findNode(target)
	if node == nil {
		return false
	}

	if flag == 0 {
		node.value = newValue
		return true
	}

	newNode := &listNode{value: newValue}
	l.length++
	// insert after
	if flag > 0 {
		next := node.next
		node.next = newNode
		newNode.prev = node
		if next == nil {
			l.tail = newNode
		} else {
			newNode.next = next
			next.prev = newNode
		}

		return true
	}

	// insert before
	prev := node.prev
	node.prev = newNode
	newNode.next = node
	if prev == nil {
		l.head = newNode
	} else {
		newNode.prev = prev
		prev.next = newNode
	}

	return true
}

func (l *List) findNode(val string) *listNode {
	var cur = l.head
	for cur != nil {
		if cur.value == val {
			return cur
		}
		cur = cur.next
	}

	return nil
}

/*
 * --- debug ---
 */

func (l *List) print() {
	var (
		sb   = new(strings.Builder)
		temp = l.head
	)

	sb.WriteString("[")
	for temp != nil {
		sb.WriteString(" " + temp.value)
		temp = temp.next
		if temp != nil {
			sb.WriteString(",")
		}
	}
	sb.WriteString(" ]")
	fmt.Printf("len:%d, values: %s\n", l.length, sb.String())
}

/*
 * Commands
 */

func loadAndCheckList(key string, checkLen bool) (*List, error) {
	info, err := loadKeyInfo(key, KeyTypeList)
	if err != nil {
		return nil, err
	}

	list := info.Value.(*List)
	if checkLen && list.length == 0 {
		return nil, ErrNil
	}

	return list, nil
}

func LPush(key string, values ...string) (ln int, err error) {
	info, err := loadKeyInfo(key, KeyTypeList)
	if err == ErrNil {
		list := newListByLPush(values...)
		defaultCache.keys.Set(key, &KeyInfo{
			Type:  KeyTypeList,
			Value: list,
		})

		return list.length, nil
	}

	if err != nil {
		return 0, err
	}

	list := info.Value.(*List)
	for _, value := range values {
		list.LPush(value)
	}

	return list.length, nil
}

func LPop(key string, count int) (values []string, err error) {
	info, err := loadKeyInfo(key, KeyTypeList)
	if err != nil {
		return nil, err
	}

	list := info.Value.(*List)
	for count > 0 && list.length > 0 {
		value, ok := list.LPop()
		if !ok {
			return
		}

		values = append(values, value)
		count--
	}

	return
}

func RPush(key string, values ...string) (ln int, err error) {
	info, err := loadKeyInfo(key, KeyTypeList)
	if err == ErrNil {
		list := newListByRPush(values...)
		defaultCache.keys.Set(key, &KeyInfo{
			Type:  KeyTypeList,
			Value: list,
		})

		return list.length, nil
	}

	if err != nil {
		return 0, err
	}

	list := info.Value.(*List)
	for _, value := range values {
		list.RPush(value)
	}

	return list.length, nil
}

func RPop(key string, count int) (values []string, err error) {
	list, err := loadAndCheckList(key, false)
	if err != nil {
		return nil, err
	}

	for count > 0 && list.length > 0 {
		value, ok := list.RPop()
		if !ok {
			return
		}

		values = append(values, value)
		count--
	}

	return
}

func LLen(key string) (ln int, err error) {
	list, err := loadAndCheckList(key, false)
	if err != nil {
		return 0, err
	}

	ln = list.length
	return
}

func LRange(key string, start, stop int) (values []string, err error) {
	list, err := loadAndCheckList(key, true)
	if err != nil {
		return nil, err
	}

	values = list.LRange(start, stop)
	return
}

// LRem Removes the first count occurrences of elements equal to element from the list stored at key
func LRem(key string, count int, value string) (n int, err error) {
	list, err := loadAndCheckList(key, true)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return list.LRemCountFromHead(value, count), nil
	}

	if count < 0 {
		return list.LRemCountFromTail(value, -count), nil
	}

	return list.LRemCountFromHead(value, list.length), nil
}

func LIndex(key string, index int) (value string, err error) {
	list, err := loadAndCheckList(key, true)
	if err != nil {
		return "", err
	}

	value, ok := list.LIndex(index)
	if !ok {
		return "", ErrNil
	}

	return value, nil
}

func LSet(key string, index int, value string) error {
	list, err := loadAndCheckList(key, true)
	if err != nil {
		return err
	}

	ok := list.LSet(index, value)
	if !ok {
		return ErrNil
	}

	return nil
}

// LInsert insert newValue into list
// if flag == 0 then replace target
// if flag > 0 insert newValue after target
// if flag < 0 insert newValue before target
// n is new length of list after insert
func LInsert(key, target, newValue string, pos int) (n int, err error) {
	list, err := loadAndCheckList(key, true)
	if err != nil {
		return 0, err
	}

	if !list.LInsert(target, newValue, pos) {
		return 0, ErrNil
	}

	return list.length, nil
}
