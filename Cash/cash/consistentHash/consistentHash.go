package consistentHash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32
type ConsistentHash struct {
	numReplicas int
	hash        Hash
	vnodes      []int
	vnode2Node  map[uint32]string
}

func Default(numReplicas int, nodes []string) *ConsistentHash {
	ch := &ConsistentHash{
		numReplicas: numReplicas,
		hash:        crc32.ChecksumIEEE,
		vnode2Node:  map[uint32]string{},
	}
	ch.init(nodes)
	return ch
}

func New(numReplicas int, nodes []string, hash Hash) *ConsistentHash {
	ch := &ConsistentHash{
		numReplicas: numReplicas,
		hash:        hash,
		vnode2Node:  map[uint32]string{},
	}
	if ch.hash == nil {
		ch.hash = crc32.ChecksumIEEE
	}
	ch.init(nodes)
	return ch
}

func (ch *ConsistentHash) init(nodes []string) {
	for _, node := range nodes {
		for i := 0; i < ch.numReplicas; i++ {
			vnodeStr := strconv.Itoa(i) + node
			hash := ch.hash([]byte(vnodeStr))
			ch.vnode2Node[hash] = node
			ch.vnodes = append(ch.vnodes, int(hash))
		}
	}
	sort.Ints(ch.vnodes)
}

func (ch *ConsistentHash) AddNode(node string) {
	for i := 0; i < ch.numReplicas; i++ {
		vnodeStr := strconv.Itoa(i) + node
		hash := ch.hash([]byte(vnodeStr))
		ch.vnode2Node[hash] = node
		ch.vnodes = append(ch.vnodes, int(hash))
	}
	sort.Ints(ch.vnodes)
}

func (ch *ConsistentHash) DeleteNode(node string) {
	i := 0
	for i < len(ch.vnodes) {
		vnode := uint32(ch.vnodes[i])
		if _node, ok := ch.vnode2Node[vnode]; ok && _node == node {
			ch.vnodes = append(ch.vnodes[:i], ch.vnodes[i+1:]...)
			delete(ch.vnode2Node, vnode)
		} else {
			i++
		}
	}
}

func (ch *ConsistentHash) Get(key string) string {
	hash := ch.hash([]byte(key))
	n := len(ch.vnodes)
	idx := sort.Search(n, func(i int) bool {
		return uint32(ch.vnodes[i]) >= hash
	})

	return ch.vnode2Node[uint32(ch.vnodes[idx%n])]
}
