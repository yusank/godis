package datastruct

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// sortedSet implement via skip list

// zSet is object contain skip list and map which store key-value pair
type zSet struct {
	m   sync.Map   // store key and value
	zsl *zSkipList // skip list
}

type zSkipList struct {
	head, tail *zSkipListNode
	length     uint // 总长度
	level      int  // 最大高度
}

type zSkipListNode struct {
	value    string
	score    float64
	backward *zSkipListNode
	levels   []*zSkipListLeve
}

type zSkipListLeve struct {
	forward *zSkipListNode
	span    uint // 当前 level 到下一个节点的跨度
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	ZSkipListMaxLevel = 1 << 5 // enough for 2^64 elements
	ZSkipListP        = 4      // 随机因子
)

func newZSkipList() *zSkipList {
	zsl := &zSkipList{
		level: 1,
		head:  newZslNode(ZSkipListMaxLevel, 0, ""),
	}

	return zsl
}

func (zsl *zSkipList) print() {
	cur := zsl.head
	fmt.Println("zset length:", zsl.length)
	for cur != nil {
		fmt.Printf("score:%-2.0f, value:%-10s, height:%-2d\n", cur.score, cur.value, len(cur.levels))
		if len(cur.levels) > 0 {
			cur = cur.levels[0].forward
		} else {
			cur = nil
		}
	}
}

func newZslNode(level int, score float64, value string) *zSkipListNode {
	node := &zSkipListNode{
		value:  value,
		score:  score,
		levels: make([]*zSkipListLeve, level),
	}

	for i := 0; i < level; i++ {
		node.levels[i] = &zSkipListLeve{}
	}

	return node
}

func zslRandomLevel() int {
	var (
		level = 1
	)
	/* 随机的同时降低更高的level 的出现,怎么做呢?
	* Redis源码: (random()&0xFFFF) < (0.25 * 0xFFFF)
	* 这里的实现如下,原理是,随机一个数,其为偶数的概率 50%,连续两次偶数概率为 25%,以此类推
	* 连续次数越高概率越低,从而保证 level 更多的分布于低位level
	 */
	for rand.Int()%2 == 0 {
		level++
	}

	if level < ZSkipListMaxLevel {
		return level
	}
	return ZSkipListMaxLevel
}

func tripleOp(cond bool, trueVal, falseVal uint) uint {
	if cond {
		return trueVal
	}

	return falseVal
}

func (zsl *zSkipList) insert(score float64, value string) *zSkipListNode {
	var (
		update   = make([]*zSkipListNode, ZSkipListMaxLevel) // 记录更新
		x        *zSkipListNode
		rank     = make([]uint, ZSkipListMaxLevel)
		i, level int
	)

	x = zsl.head
	for i = zsl.level - 1; i >= 0; i-- {
		rank[i] = tripleOp(i == zsl.level-1, 0, rank[i+1])
		for x.levels[i].forward != nil &&
			(x.levels[i].forward.score < score ||
				(x.levels[i].forward.score == score &&
					(x.levels[i].forward.value < value))) {
			rank[i] += x.levels[i].span
			x = x.levels[i].forward
		}

		update[i] = x
	}

	level = zslRandomLevel()
	if level > zsl.level {
		for i = zsl.level; i < level; i++ {
			rank[i] = 0
			update[i] = zsl.head
			update[i].levels[i].span = zsl.length
		}

		zsl.level = level
	}

	x = newZslNode(level, score, value)
	for i = 0; i < level; i++ {
		x.levels[i].forward = update[i].levels[i].forward
		update[i].levels[i].forward = x

		x.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i])
		update[i].levels[i].span = (rank[0] - rank[i]) + 1
	}

	for i = level; i < zsl.level; i++ {
		update[i].levels[i].span++
	}

	if update[0] != zsl.head {
		x.backward = update[0]
	}
	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x
	} else {
		zsl.tail = x
	}

	zsl.length++

	return x
}

// return 1 if found and update element,otherwise return 0
func (zsl *zSkipList) updateScore(curScore, newScore float64, value string) *zSkipListNode {
	var update = make([]*zSkipListNode, ZSkipListMaxLevel)

	// find where the match element
	x := zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil &&
			(x.levels[i].forward.score < curScore ||
				(x.levels[i].forward.score == curScore &&
					x.levels[i].forward.value < value)) {
			x = x.levels[i].forward
		}
		update[i] = x
	}

	x = x.levels[0].forward
	if x == nil || x.score != curScore || x.value != value {
		return nil
	}

	// if after update score, no need to change node position
	// then update and return
	if (x.backward == nil || x.backward.score < newScore) &&
		(x.levels[0].forward == nil || x.levels[0].forward.score > newScore) {
		x.score = newScore
		return x
	}

	// del and insert
	zsl.deleteNode(x, update)
	newNode := zsl.insert(newScore, value)
	return newNode
}

func (zsl *zSkipList) deleteNode(node *zSkipListNode, update []*zSkipListNode) {
	for i := 0; i < zsl.level; i++ {
		if update[i].levels[i].forward == node {
			// 如果当前位置&&当前 level 的下一个是要删除的 node
			// 当前节点 level 的宽度 加上被删除节点的跨度-1
			// 当前节点 level 的下一个指向被删除的 node 的下一个
			update[i].levels[i].span += node.levels[i].span - 1
			update[i].levels[i].forward = node.levels[i].forward
		} else {
			// 否则 跨度减一
			update[i].levels[i].span -= 1
		}
	}

	// 被删除节点的下一个节点的上一个节点指向被删除节点的上一个节点
	if node.levels[0].forward != nil {
		node.levels[0].forward.backward = node.backward
	}

	// level 可能需要减少
	for zsl.level > 1 && zsl.head.levels[zsl.level-1].forward == nil {
		zsl.level--
	}

	zsl.length--
}

// delete an element with matching by score and value
//  if node is not nil then assign the deleted element to node
func (zsl *zSkipList) delete(score float64, value string, node **zSkipListNode) int {
	// match
	var update = make([]*zSkipListNode, ZSkipListMaxLevel)
	x := zsl.head

	for i := zsl.level - 1; i >= 0; i-- {
		if x.levels[i].forward != nil &&
			(x.levels[i].forward.score < score ||
				(x.levels[i].forward.score == score &&
					x.levels[i].forward.value < value)) {
			x = x.levels[i].forward
		}

		update[i] = x
	}

	// x is the element what we seek for
	x = x.levels[0].forward
	if x != nil && x.score == score && x.value == value {
		zsl.deleteNode(x, update)
		if node != nil {
			*node = x
		}

		return 1
	}

	return 0
}

func (zs *zSet) zAdd(score float64, value string, flag int) int {
	de := zs.findFromMap(value)
	if de == nil {
		node := zs.zsl.insert(score, value)
		zs.m.Store(value, withValue(&node.score))
		return 1
	}

	// nx flag
	if flag&ZAddInNx != 0 {
		// exist
		return 0
	}

	// exists
	oldScore := *(de.getValue().(*float64))
	if flag&ZAddInIncr != 0 {
		score += oldScore
	}

	if score != oldScore {
		node := zs.zsl.updateScore(oldScore, score, value)
		if node == nil {
			return 0
		}

		// update score
		de.setValue(&node.score)
		return 1
	}

	return 0
}

func (zs *zSet) findFromMap(key string) *dictEntry {
	v, ok := zs.m.Load(key)
	if !ok {
		return nil
	}

	return v.(*dictEntry)
}

/*
 * Commands
 */

func loadAndCheckZSet(key string, checkLen bool) (*zSet, error) {
	info, err := loadKeyInfo(key, KeyTypeSortedSet)
	if err != nil {
		return nil, err
	}

	zs := info.Value.(*zSet)
	if checkLen && zs.zsl.length == 0 {
		return nil, ErrNil
	}

	return zs, nil
}

type ZSetMember struct {
	Score float64
	Value string
}

func ZAdd(key string, members []*ZSetMember, flag int) (int, error) {
	zs, err := loadAndCheckZSet(key, false)
	if err != nil {
		return 0, err
	}

	for _, m := range members {
		zs.zAdd(m.Score, m.Value, flag)
	}

	return 0, nil
}
