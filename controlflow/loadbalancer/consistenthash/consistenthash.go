package consistenthash

import (
	"errors"
	"github.com/spaolacci/murmur3"
	"sort"
	"strconv"
	"sync"
)

// ErrEmpty ...
var ErrEmpty = errors.New("empty hash ring")

// DefaultVirtual ...
const DefaultVirtual = 32

// DefaultHashFunc ...
var DefaultHashFunc = hash

// ConsistentHash struct
type ConsistentHash struct {
	sync.RWMutex // concurrent

	ring     map[uint32]string   // hash to node
	nodesSet map[string]struct{} // set of nodes
	virtual  int                 // physical node to several virtual nodes
	cnt      int                 // nodes number

	sorted []uint32 // sorted hash node

	hash HashFunc
}

// HashFunc ...
type HashFunc func([]byte) uint32

// Option ...
type Option struct {
	virtual int
	hash    HashFunc
}

var defaultOption = Option{
	virtual: DefaultVirtual,
	hash:    DefaultHashFunc,
}

func hash(data []byte) uint32 {
	m := murmur3.New32()
	_, _ = m.Write(data)
	return m.Sum32()
}

func virtualKey(node string, idx int) string {
	return node + "#" + strconv.Itoa(idx)
}

// NewCH creates a new ConsistentHash
func NewCH(option *Option) *ConsistentHash {
	if option == nil {
		option = &defaultOption
	}
	return &ConsistentHash{
		RWMutex:  sync.RWMutex{},
		ring:     make(map[uint32]string),
		nodesSet: make(map[string]struct{}),
		virtual:  option.virtual,
		cnt:      0,
		hash:     option.hash,
	}
}

// AddNode adds one node to ConsistentHash
func (ch *ConsistentHash) AddNode(node string) {
	ch.Lock()
	defer ch.Unlock()
	// already in
	if _, found := ch.nodesSet[node]; found {
		return
	}
	for i := 0; i < ch.virtual; i++ {
		k := ch.hash([]byte(virtualKey(node, i)))
		ch.ring[k] = node
		ch.sorted = append(ch.sorted, k)
	}
	// update
	ch.nodesSet[node] = struct{}{}
	ch.cnt++
	sortHashes(ch.sorted)
	return
}

// AddNodes add nodes to ConsistentHash
func (ch *ConsistentHash) AddNodes(nodes []string) {
	ch.Lock()
	defer ch.Unlock()
	for _, node := range nodes {
		// already in
		if _, found := ch.nodesSet[node]; found {
			continue
		}
		for i := 0; i < ch.virtual; i++ {
			k := ch.hash([]byte(virtualKey(node, i)))
			ch.ring[k] = node
			ch.sorted = append(ch.sorted, k)
		}
		// update
		ch.nodesSet[node] = struct{}{}
		ch.cnt++
	}
	sortHashes(ch.sorted)
}

// DeleteNode deletes one node from ConsistentHash
func (ch *ConsistentHash) DeleteNode(node string) {
	ch.Lock()
	defer ch.Unlock()
	// not in
	if _, found := ch.nodesSet[node]; !found {
		return
	}
	for i := 0; i < ch.virtual; i++ {
		k := ch.hash([]byte(virtualKey(node, i)))
		delete(ch.ring, k)
	}
	// update
	delete(ch.nodesSet, node)
	ch.cnt--
	ch.sorted = ch.sorted[:0]
	for k := range ch.ring {
		ch.sorted = append(ch.sorted, k)
	}
	sortHashes(ch.sorted)
}

// DeleteNodes deletes nodes from ConsistentHash
func (ch *ConsistentHash) DeleteNodes(nodes []string) {
	ch.Lock()
	defer ch.Unlock()
	for _, node := range nodes {
		// not in
		if _, found := ch.nodesSet[node]; !found {
			continue
		}
		for i := 0; i < ch.virtual; i++ {
			k := ch.hash([]byte(virtualKey(node, i)))
			delete(ch.ring, k)
		}
		// update
		delete(ch.nodesSet, node)
		ch.cnt--
	}
	ch.sorted = ch.sorted[:0]
	for k := range ch.ring {
		ch.sorted = append(ch.sorted, k)
	}
	sortHashes(ch.sorted)
}

// Get gets one node from ConsistentHash
func (ch *ConsistentHash) Get(key string) (node string, err error) {
	ch.RLock()
	defer ch.RUnlock()
	if ch.cnt == 0 {
		return "", ErrEmpty
	}
	hash := ch.hash([]byte(key))
	// len(ch.sorted) == ch.cnt * ch.virtual
	l := len(ch.sorted)
	idx := sort.Search(l, func(i int) bool {
		return ch.sorted[i] >= hash
	})
	// ring end -> ring start
	if idx == l {
		idx = 0
	}
	return ch.ring[ch.sorted[idx]], nil
}

// Set sets nodes of ConsistentHash, it will replace former nodes with input ones
func (ch *ConsistentHash) Set(nodes []string) {
	ch.Lock()
	defer ch.Unlock()
	markDeleted := make([]string, 0, ch.cnt)
	newNodes := make(map[string]struct{}, len(nodes))
	for _, node := range nodes {
		newNodes[node] = struct{}{}
	}
	// delete
	for node := range ch.nodesSet {
		if _, found := newNodes[node]; !found {
			markDeleted = append(markDeleted, node)
		}
	}
	for _, n := range markDeleted {
		for i := 0; i < ch.virtual; i++ {
			k := ch.hash([]byte(virtualKey(n, i)))
			delete(ch.ring, k)
		}
		// update
		delete(ch.nodesSet, n)
		ch.cnt--
	}

	// add
	for n := range newNodes {
		if _, found := ch.nodesSet[n]; !found {
			for i := 0; i < ch.virtual; i++ {
				k := ch.hash([]byte(virtualKey(n, i)))
				ch.ring[k] = n
			}
			// update
			ch.nodesSet[n] = struct{}{}
			ch.cnt++
		}
	}

	// update
	ch.sorted = ch.sorted[:0]
	for k := range ch.ring {
		ch.sorted = append(ch.sorted, k)
	}
	sortHashes(ch.sorted)
	return
}

// Nodes returns nodes of ConsistentHash as a slice of string
func (ch *ConsistentHash) Nodes() []string {
	ch.RLock()
	defer ch.RUnlock()
	ret := make([]string, 0, ch.cnt)
	for k := range ch.nodesSet {
		ret = append(ret, k)
	}
	return ret
}

func sortHashes(hashes []uint32) {
	sort.Slice(hashes, func(i, j int) bool {
		return hashes[i] < hashes[j]
	})
}
