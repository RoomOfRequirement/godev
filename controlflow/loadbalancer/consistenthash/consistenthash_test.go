package consistenthash

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"hash/crc32"
	"hash/fnv"
	"math/rand"
	"testing"
)

func TestNewCH(t *testing.T) {
	ch := NewCH(nil)
	assert.Equal(t, DefaultVirtual, ch.virtual)
	assert.Equal(t, 0, ch.cnt)
	assert.Equal(t, 0, len(ch.nodesSet))
	assert.Equal(t, 0, len(ch.ring))
	assert.Equal(t, 0, len(ch.sorted))

	ch = NewCH(&Option{virtual: 64, hash: DefaultHashFunc})
	assert.Equal(t, 64, ch.virtual)
}

func TestConsistentHash_AddNode(t *testing.T) {
	ch := NewCH(nil)
	ch.AddNode("0.0.0.0")
	assert.Equal(t, DefaultVirtual, ch.virtual)
	assert.Equal(t, 1, ch.cnt)
	assert.Equal(t, 1, len(ch.nodesSet))
	assert.Equal(t, ch.virtual, len(ch.ring))
	assert.Equal(t, ch.virtual, len(ch.sorted))

	ch.AddNode("0.0.0.0")
	assert.Equal(t, DefaultVirtual, ch.virtual)
	assert.Equal(t, 1, ch.cnt)
	assert.Equal(t, 1, len(ch.nodesSet))
	assert.Equal(t, ch.virtual, len(ch.ring))
	assert.Equal(t, ch.virtual, len(ch.sorted))

	ch.AddNode("0.0.0.1")
	assert.Equal(t, DefaultVirtual, ch.virtual)
	assert.Equal(t, 2, ch.cnt)
	assert.Equal(t, 2, len(ch.nodesSet))
	assert.Equal(t, ch.virtual*2, len(ch.ring))
	assert.Equal(t, ch.virtual*2, len(ch.sorted))
}

func TestConsistentHash_AddNodes(t *testing.T) {
	ch := NewCH(nil)
	ch.AddNodes([]string{"0.0.0.0", "0.0.0.0", "0.0.0.1", "0.0.0.2", "0.0.0.2"})
	assert.Equal(t, DefaultVirtual, ch.virtual)
	assert.Equal(t, 3, ch.cnt)
	assert.Equal(t, 3, len(ch.nodesSet))
	assert.Equal(t, ch.virtual*3, len(ch.ring))
	assert.Equal(t, ch.virtual*3, len(ch.sorted))
}

func TestConsistentHash_DeleteNode(t *testing.T) {
	ch := NewCH(nil)
	ch.AddNode("0.0.0.0")
	ch.DeleteNode("0.0.0.1")
	assert.Equal(t, DefaultVirtual, ch.virtual)
	assert.Equal(t, 1, ch.cnt)
	assert.Equal(t, 1, len(ch.nodesSet))
	assert.Equal(t, ch.virtual, len(ch.ring))
	assert.Equal(t, ch.virtual, len(ch.sorted))

	ch.DeleteNode("0.0.0.0")
	assert.Equal(t, DefaultVirtual, ch.virtual)
	assert.Equal(t, 0, ch.cnt)
	assert.Equal(t, 0, len(ch.nodesSet))
	assert.Equal(t, 0, len(ch.ring))
	assert.Equal(t, 0, len(ch.sorted))

	ch.AddNodes([]string{"0.0.0.0", "0.0.0.0", "0.0.0.1", "0.0.0.2", "0.0.0.2"})
	ch.DeleteNode("0.0.0.0")
	assert.Equal(t, DefaultVirtual, ch.virtual)
	assert.Equal(t, 2, ch.cnt)
	assert.Equal(t, 2, len(ch.nodesSet))
	assert.Equal(t, ch.virtual*2, len(ch.ring))
	assert.Equal(t, ch.virtual*2, len(ch.sorted))
}

func TestConsistentHash_DeleteNodes(t *testing.T) {
	ch := NewCH(nil)
	ch.AddNodes([]string{"0.0.0.0", "0.0.0.0", "0.0.0.1", "0.0.0.2", "0.0.0.2"})
	ch.DeleteNodes([]string{"0.0.0.0", "0.0.0.3"})
	assert.Equal(t, DefaultVirtual, ch.virtual)
	assert.Equal(t, 2, ch.cnt)
	assert.Equal(t, 2, len(ch.nodesSet))
	assert.Equal(t, ch.virtual*2, len(ch.ring))
	assert.Equal(t, ch.virtual*2, len(ch.sorted))
}

func TestConsistentHash_Get(t *testing.T) {
	ch := NewCH(nil)
	node, err := ch.Get("hello")
	assert.Error(t, err)
	assert.Equal(t, ErrEmpty, err)
	assert.Empty(t, node)
	ch.AddNodes([]string{"0.0.0.0", "0.0.0.0", "0.0.0.1", "0.0.0.2", "0.0.0.2"})
	node, err = ch.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0.0", node)
}

func TestConsistentHash_Set(t *testing.T) {
	ch := NewCH(nil)
	ch.AddNodes([]string{"0.0.0.0", "0.0.0.0", "0.0.0.1", "0.0.0.2", "0.0.0.2"})

	ch.Set([]string{"0.0.0.3", "0.0.0.1", "0.0.0.6"})
	assert.Equal(t, DefaultVirtual, ch.virtual)
	assert.Equal(t, 3, ch.cnt)
	assert.Equal(t, 3, len(ch.nodesSet))
	assert.Equal(t, ch.virtual*3, len(ch.ring))
	assert.Equal(t, ch.virtual*3, len(ch.sorted))
	assert.EqualValues(t, map[string]struct{}{
		"0.0.0.1": {},
		"0.0.0.3": {},
		"0.0.0.6": {},
	}, ch.nodesSet)
}

func TestConsistentHash_Nodes(t *testing.T) {
	ch := NewCH(nil)
	ch.AddNodes([]string{"0.0.0.0", "0.0.0.0", "0.0.0.1", "0.0.0.2", "0.0.0.2"})
	expected := map[string]struct{}{
		"0.0.0.0": {},
		"0.0.0.1": {},
		"0.0.0.2": {},
	}
	nodes := ch.Nodes()
	assert.Equal(t, len(expected), len(nodes))
	for _, node := range nodes {
		if _, found := expected[node]; !found {
			t.Fail()
		}
	}
}

func TestDistribution(t *testing.T) {
	num := 100000
	dist := testDistribution(t, hash, num)
	for k, v := range dist {
		t.Logf("%s: %0.3f", "(murmur3)"+k, float64(v)/float64(num))
	}
	dist = testDistribution(t, crc32.ChecksumIEEE, num)
	for k, v := range dist {
		t.Logf("%s: %0.3f", "(crc32)"+k, float64(v)/float64(num))
	}
	dist = testDistribution(t, func(bytes []byte) uint32 {
		h := fnv.New32a()
		_, _ = h.Write(bytes)
		return h.Sum32()
	}, num)
	for k, v := range dist {
		t.Logf("%s: %0.3f", "(fnv32)"+k, float64(v)/float64(num))
	}
}

/*
 * (murmur3)0.0.0.0: 0.278
 * (murmur3)0.0.0.1: 0.448
 * (murmur3)0.0.0.2: 0.274
 * (crc32)0.0.0.0: 0.298
 * (crc32)0.0.0.1: 0.300
 * (crc32)0.0.0.2: 0.401
 * (fnv32)0.0.0.0: 0.639
 * (fnv32)0.0.0.1: 0.196
 * (fnv32)0.0.0.2: 0.165
 */

// TODO: some other metrics like stderr?
func TestRehash(t *testing.T) {
	num := 100000
	delta := testRehash(t, hash, num)
	t.Logf("%s: %0.3f", "murmur3 rehash rate", float64(delta)/float64(num))
	delta = testRehash(t, crc32.ChecksumIEEE, num)
	t.Logf("%s: %0.3f", "crc32 rehash rate", float64(delta)/float64(num))
	delta = testRehash(t, func(bytes []byte) uint32 {
		h := fnv.New32a()
		_, _ = h.Write(bytes)
		return h.Sum32()
	}, num)
	t.Logf("%s: %0.3f", "fnv32 rehash rate", float64(delta)/float64(num))
}

/*
 * murmur3 rehash rate: 0.241
 * crc32 rehash rate: 0.199
 * fnv32 rehash rate: 0.156
 */

// add one node
func testRehash(t *testing.T, hashFunc HashFunc, num int) int {
	ch := NewCH(&Option{virtual: DefaultVirtual, hash: hashFunc})
	ch.AddNodes([]string{"0.0.0.0", "0.0.0.1", "0.0.0.2"})
	dist := make(map[string]int)
	buf := make([]byte, 12)
	for i := 0; i < num; i++ {
		_, err := rand.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		r, err := ch.Get(base64.StdEncoding.EncodeToString(buf))
		if err != nil {
			t.Fatal(err)
		}
		dist[r]++
	}
	ch.AddNode("0.0.0.3")
	newDist := make(map[string]int)
	for i := 0; i < num; i++ {
		_, err := rand.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		r, err := ch.Get(base64.StdEncoding.EncodeToString(buf))
		if err != nil {
			t.Fatal(err)
		}
		newDist[r]++
	}
	ret := 0
	for k, v := range dist {
		ret += v - newDist[k]
	}
	return ret
}

func testDistribution(t *testing.T, hashFunc HashFunc, num int) map[string]int {
	ch := NewCH(&Option{virtual: DefaultVirtual, hash: hashFunc})
	ch.AddNodes([]string{"0.0.0.0", "0.0.0.1", "0.0.0.2"})
	dist := make(map[string]int)
	buf := make([]byte, 12)
	for i := 0; i < num; i++ {
		_, err := rand.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		r, err := ch.Get(base64.StdEncoding.EncodeToString(buf))
		if err != nil {
			t.Fatal(err)
		}
		dist[r]++
	}
	return dist
}
