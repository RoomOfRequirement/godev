package consistenthash

/*
 * From the paper "A Fast, Minimal Memory, Consistent Hash Algorithm" by John Lamping, Eric Veach (Google, 2014).
 * http://arxiv.org/abs/1406.2294
 */

// JumpHash ...
//	evenly distributed: 3 nodes, every node will have 1/3 keys
//	if key is integer, no need to hash it
//	faster than `Get` implementation in `consistenthash.go`
//	it can evenly distribute keys when numBuckets changes
//	TODO: implement consistent hash based on JumpHash?
func JumpHash(key uint64, numBuckets int) int32 {
	var b int64 = -1
	var j int64 = 0
	for j < int64(numBuckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64((b+1)*(int64(1)<<31)) / float64((key>>33)+1))
	}
	return int32(b)
}
