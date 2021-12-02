package datastruct

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

// sortedSet implement via skip list

// zSet is object contain skip list and map which store key-value pair
type zSet struct {
	m sync.Map // store key and value
	// 在元素少于 100 & 每个元素大小小于 64 的时候,Redis 实际上用的是 zipList 这里作为知识点提了一下
	// 	除非遇到性能问题,否则不准备同时支持 zipList 和 skipList
	zsl *zSkipList // skip list
}

func newZSet() *zSet {
	return &zSet{
		m:   sync.Map{},
		zsl: newZSkipList(),
	}
}

type dictEntry struct {
	value interface{}
}

func (de *dictEntry) getValue() interface{} {
	return de.value
}

func (de *dictEntry) setValue(v interface{}) {
	de.value = v
}

func withValue(v interface{}) *dictEntry {
	return &dictEntry{value: v}
}

type zSkipList struct {
	head, tail *zSkipListNode
	length     int // 总长度
	level      int // 最大高度
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
			update[i].levels[i].span = uint(zsl.length)
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

// make sure there match an element in the skip list
func (zsl *zSkipList) rank(score float64, value string) uint {
	var rank uint
	x := zsl.head
LevelLoop:
	for i := zsl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil &&
			(x.levels[i].forward.score < score ||
				(x.levels[i].forward.score == score &&
					(x.levels[i].forward.value <= value))) {
			rank += x.levels[i].span
			x = x.levels[i].forward
			if x.score == score && x.value == value {
				break LevelLoop
			}
		}
	}

	return rank
}

// minOpen == true => (min,  for example (1,true,2,false) => zscore xxx (1 2
func (zsl *zSkipList) count(min float64, minOpen bool, max float64, maxOpen bool) int {
	// -inf +inf
	if math.IsInf(min, -1) && math.IsInf(max, 1) {
		return zsl.length
	}

	minRank := zsl.minScoreRank(min, minOpen)
	maxRank := zsl.maxScoreRank(max, maxOpen)
	fmt.Println(minRank, maxRank)

	// -inf:max
	if minRank <= 0 {
		return maxRank
	}

	if maxRank < 0 {
		return zsl.length - minRank + 1
	}

	return maxRank - minRank + 1
}

// return -1 if score is inf
func (zsl *zSkipList) minScoreRank(score float64, openInterval bool) int {
	if math.IsInf(score, 0) {
		return -1
	}

	var (
		rank uint
		x    = zsl.head
		cur  = zsl.head
	)
loop:
	for i := zsl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil &&
			(x.levels[i].forward.score <= score) {
			rank += x.levels[i].span
			cur = x
			x = x.levels[i].forward

			// got the position
			if x.score > score || (openInterval && x.score == score) {
				// score not exit in the list, it just between two elements
				if cur.score < score && x.score > score {
					rank += 1
				}

				if x.score == score && openInterval {
					rank += 1
				}

				break loop
			}
		}
	}

	return int(rank)
}

// return -1 if score is inf
func (zsl *zSkipList) maxScoreRank(score float64, openInterval bool) int {
	if math.IsInf(score, 0) {
		return -1
	}

	var (
		rank uint
		x    = zsl.head
	)
loop:
	for i := zsl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil &&
			(x.levels[i].forward.score <= score) {
			rank += x.levels[i].span
			x = x.levels[i].forward

			// got the position
			if x.score > score || (openInterval && x.score == score) {
				if x.score == score && openInterval {
					rank -= 1
				}

				break loop
			}
		}
	}

	return int(rank)
}

func (zsl *zSkipList) zRange(start, stop int, withScores bool) []string {
	if start < 0 {
		start = start + zsl.length
		if start < 0 {
			start = 0
		}
	}

	if stop < 0 {
		stop = stop + zsl.length
	}

	if start > stop || start >= zsl.length {
		return nil
	}
	if stop >= zsl.length {
		stop = zsl.length - 1
	}

	node := zsl.findElementByRank(uint(start) + 1)
	var (
		rangeLen = stop - start + 1
		result   []string
	)
	for rangeLen > 0 {
		result = append(result, node.value)
		if withScores {
			result = append(result, strconv.FormatFloat(node.score, 'g', -1, 64))
		}

		node = node.levels[0].forward
		rangeLen--
	}

	return result
}

func (zsl *zSkipList) zRangeByScore(min float64, minOpen bool, max float64, maxOpen bool, withScores bool) []string {
	// -inf +inf
	if math.IsInf(min, -1) && math.IsInf(max, 1) {
		return zsl.zRange(0, zsl.length, withScores)
	}

	minRank := zsl.minScoreRank(min, minOpen)
	maxRank := zsl.maxScoreRank(max, maxOpen)
	fmt.Println(minRank, maxRank)
	// -inf:max
	if minRank <= 0 {
		return zsl.zRange(0, maxRank-1, withScores)
	}

	if maxRank < 0 {
		return zsl.zRange(minRank-1, zsl.length, withScores)
	}

	return zsl.zRange(minRank-1, maxRank-1, withScores)
}

func (zsl *zSkipList) findElementByRank(rank uint) *zSkipListNode {
	var (
		x         = zsl.head
		traversed uint
	)

	for i := zsl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && traversed+x.levels[i].span <= rank {
			traversed += x.levels[i].span
			x = x.levels[i].forward
		}
		if traversed == rank {
			return x
		}
	}

	return nil
}

func (zs *zSet) zAdd(score float64, value string, flag int) *zSkipListNode {
	de := zs.loadDictEntry(value)
	if de == nil {
		// xx flag
		if flag&ZAddInXx != 0 {
			return nil
		}

		node := zs.zsl.insert(score, value)
		zs.m.Store(value, withValue(&node.score))
		return node
	}

	// nx flag
	if flag&ZAddInNx != 0 {
		// exist
		return nil
	}

	// exists
	oldScore := *(de.getValue().(*float64))
	if flag&ZAddInIncr != 0 {
		score += oldScore
	}

	if score != oldScore {
		node := zs.zsl.updateScore(oldScore, score, value)
		if node == nil {
			return nil
		}

		// update score
		de.setValue(&node.score)
		return node
	}

	return nil
}

func (zs *zSet) loadDictEntry(key string) *dictEntry {
	v, ok := zs.m.Load(key)
	if !ok {
		return nil
	}

	return v.(*dictEntry)
}

func (zs *zSet) loadAndDeleteDictEntry(key string) *dictEntry {
	v, ok := zs.m.LoadAndDelete(key)
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
	if err != nil && err != ErrNil {
		return 0, err
	}

	if zs == nil {
		zs = newZSet()
		defaultCache.keys.Store(key, &KeyInfo{
			Type:  KeyTypeSortedSet,
			Value: zs,
		})
	}

	var cnt int
	for _, m := range members {
		if node := zs.zAdd(m.Score, m.Value, flag); node != nil {
			cnt++
		}
	}

	return cnt, nil
}

func ZScore(key, value string) (float64, error) {
	zs, err := loadAndCheckZSet(key, true)
	if err != nil {
		return 0, err
	}

	de := zs.loadDictEntry(value)
	if de == nil {
		return 0, ErrNil
	}

	return *(de.value.(*float64)), nil
}

func ZRank(key, value string) (uint, error) {
	zs, err := loadAndCheckZSet(key, true)
	if err != nil {
		return 0, err
	}

	de := zs.loadDictEntry(value)
	if de == nil {
		return 0, ErrNil
	}

	score := *(de.value.(*float64))
	return zs.zsl.rank(score, value), nil
}

func ZRem(key string, values ...string) (int, error) {
	zs, err := loadAndCheckZSet(key, true)
	if err != nil {
		return 0, err
	}

	var cnt int
	for _, value := range values {
		de := zs.loadAndDeleteDictEntry(value)
		if de == nil {
			continue
		}

		score := *(de.value.(*float64))
		zs.zsl.delete(score, value, nil)
		cnt++
	}

	return cnt, nil
}

func ZCard(key string) (int, error) {
	zs, err := loadAndCheckZSet(key, true)
	if err != nil {
		return 0, err
	}

	return zs.zsl.length, nil
}

func ZCount(key, minStr, maxStr string) (int, error) {
	zs, err := loadAndCheckZSet(key, true)
	if err != nil {
		return 0, err
	}

	minScore, minOpen, err := handleFloatScoreStr(minStr)
	if err != nil {
		return 0, err
	}

	maxScore, maxOpen, err := handleFloatScoreStr(maxStr)
	if err != nil {
		return 0, err
	}

	return zs.zsl.count(minScore, minOpen, maxScore, maxOpen), nil
}

func ZIncr(key string, score float64, value string) (float64, error) {
	zs, err := loadAndCheckZSet(key, false)
	if err != nil && err != ErrNil {
		return 0, err
	}

	if zs == nil {
		zs = newZSet()
		defaultCache.keys.Store(key, &KeyInfo{
			Type:  KeyTypeSortedSet,
			Value: zs,
		})
	}

	node := zs.zAdd(score, value, ZAddInIncr)
	if node == nil {
		return 0, ErrNil
	}

	return node.score, nil
}

// ZRange not support limit for now
// todo support by lex flag
func ZRange(key string, minStr, maxStr string, flag int) ([]string, error) {
	zs, err := loadAndCheckZSet(key, true)
	if err != nil {
		return nil, err
	}

	var withScores = flag&ZRangeInWithScores != 0

	if flag&ZRangeInByScore != 0 {
		// byScore flag
		minScore, minOpen, err1 := handleFloatScoreStr(minStr)
		if err1 != nil {
			return nil, err1
		}

		maxScore, maxOpen, err1 := handleFloatScoreStr(maxStr)
		if err1 != nil {
			return nil, err1
		}

		return zs.zsl.zRangeByScore(minScore, minOpen, maxScore, maxOpen, withScores), nil
	}

	start, err := strconv.Atoi(minStr)
	if err != nil {
		return nil, ErrNotInteger
	}

	stop, err := strconv.Atoi(maxStr)
	if err != nil {
		return nil, ErrNotInteger
	}

	return zs.zsl.zRange(start, stop, withScores), nil
}

// parse score input and return value and true if is open interval
// example '(5' => 5, true,   '3' => 3, false  '-inf' => math.Inf(-1), false
func handleFloatScoreStr(str string) (float64, bool, error) {
	str = strings.ToLower(str)
	switch str {
	case "-inf":
		return math.Inf(-1), false, nil
	case "+inf":
		return math.Inf(1), false, nil
	default:
		var openInterval bool
		tmp := strings.TrimPrefix(str, "(")
		if tmp != str {
			openInterval = true
		}

		f, err := strconv.ParseFloat(tmp, 64)
		if err != nil {
			return 0, false, ErrNotFloat
		}

		return f, openInterval, nil
	}
}
