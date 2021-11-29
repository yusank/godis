package datastruct

import (
	"fmt"
	"math/rand"
	"time"
)

// sortedSet implement via skip list

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
		n     = 1 << 16
	)
	/* 这里从 Redis 源码借鉴过来
	* 源码: (random()&0xFFFF) < (ZSKIPLIST_P * 0xFFFF) // ZSKIPLIST_P=0.25
	* 通过该算法,返回的随机 level ,更接近于较低的数值,通过连续运行 10w 次后每个 level 返回的次数
	* 0 0
	* 1 50251
	* 2 24927
	* 3 12424
	* 4 6129
	* 5 3171
	* 6 1588
	* 7 774
	* 8 388
	* 9 181
	* 10 83
	* 11 41
	* 12 15
	* 13 10
	* 14 10
	* 15 5
	* 16 1
	 */
	for rand.Int()&n < n/ZSkipListP {
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
func (zsl *zSkipList) updateScore(curScore, newScore float64, value string) int {
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
		return 0
	}

	// if after update score, no need to change node position
	// then update and return
	if (x.backward == nil || x.backward.score < newScore) &&
		(x.levels[0].forward == nil || x.levels[0].forward.score > newScore) {
		x.score = newScore
		return 1
	}

	// del and insert
	zsl.deleteNode(x, update)
	zsl.insert(newScore, value)
	// todo free x
	return 1
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
		if node == nil {
			// todo should free x
		} else {
			*node = x
		}

		return 1
	}

	return 0
}

/*
 * Commands
 */

func loadAndCheckZsl(key string, checkLen bool) (*zSkipList, error) {
	info, err := loadKeyInfo(key, KeyTypeSortedSet)
	if err != nil {
		return nil, err
	}

	zsl := info.Value.(*zSkipList)
	if checkLen && zsl.length == 0 {
		return nil, ErrNil
	}

	return zsl, nil
}

type ZSetMember struct {
	Score float64
	Value string
}

func ZAdd(key string, members []*ZSetMember, options ...string) (int, error) {
	return 0, nil
}
