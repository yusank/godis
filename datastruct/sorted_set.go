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
	span    uint
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
		fmt.Printf("score:%-2.0f, value:%-10s, height:%-2d\n", cur.score, cur.value, len(cur.levels))
		if len(cur.levels) > 0 {
			cur = cur.levels[0].forward
		} else {
			cur = nil
		}
	}
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
