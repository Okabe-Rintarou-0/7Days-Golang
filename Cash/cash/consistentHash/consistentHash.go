package consistentHash

type Hash func([]byte) uint64
type ConsistentHash struct {
	numReplicas int
	hash        Hash
	vnodes      []int
	vnode2Node  map[uint64]string
}

//func New(numReplicas int, nodes []string) *ConsistentHash {
//	ch := &ConsistentHash{
//		numReplicas: numReplicas,
//	}
//	return ch
//}

func (ch *ConsistentHash) init(nodes []string) {

}
