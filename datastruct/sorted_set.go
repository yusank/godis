package datastruct

import (
	"log"
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
	span    uint
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	ZSkipListMaxLevel = 1 << 5
)

// todo when create skip list should init head of list with full level

func newZslNode(level int, score float64, value string) *zSkipListNode {
	return &zSkipListNode{
		value:  value,
		score:  score,
		levels: make([]*zSkipListLeve, level),
	}
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

func (zsl *zSkipList) print() {
	cur := zsl.head
	for cur != nil {
		log.Println(cur.score, cur.value, len(cur.levels))
	}
}

func zslRandomLevel() int {
	return rand.Intn(ZSkipListMaxLevel) + 1
}

func tripleOp(cond bool, trueVal, falseVal uint) uint {
	if cond {
		return trueVal
	}

	return falseVal
}
