package strict_quorum

import (
	"crypto/md5"
	"fmt"
	"sort"
	"unsafe"
)

func ByteArrayToInt(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

func GetHash(id string) int64 {
	data := []byte(id)
	hash := md5.Sum(data)
	return ByteArrayToInt(hash[:])
}

func Sort(int64Values []int64) []int64 {
	out := make([]int64, len(int64Values))
	int64AsIntValues := make([]int, len(int64Values))

	for i, val := range int64Values {
		int64AsIntValues[i] = int(val)
	}

	sort.Ints(int64AsIntValues)

	for i, val := range int64AsIntValues {
		out[i] = int64(val)
	}

	return out
}

func Int64ToInt(int64Values []int64) []int {
	int64AsIntValues := make([]int, len(int64Values))

	for i, val := range int64Values {
		int64AsIntValues[i] = int(val)
	}
	return int64AsIntValues
}

func (r *Ring) Search(hash int64) int {
	// index := sort.Search(len(r.NodeHashes), func(idx int) bool {
	// 	fmt.Println(r.NodeHashes[idx] >= hash)
	// 	fmt.Printf("%d >= %d\n", idx, hash)
	// 	return r.NodeHashes[idx] >= hash
	// })
	index := 0
	for idx, nodeHash := range r.NodeHashes {
		if hash <= nodeHash {
			index = idx
			break
		}
	}

	return index
}

func (r *Ring) GetNode(partitionKey string) *Node {
	hash := GetHash(partitionKey)
	index := r.Search(hash)

	nodeHash := r.NodeHashes[index]
	// return *r.Nodes[nodeHash]
	return r.NodeMap[nodeHash]
}

func (r *Ring) Replicate(partitionKey string) []*Node {
	nodesToReplicateTo := []*Node{}
	hash := GetHash(partitionKey)
	index := r.Search(hash)
	fmt.Printf("Replicating from node with hash %d\n", hash)

	// replicated nodes
	fmt.Println("Nodes to replicate to:")
	for i := 1; i < r.ReplicationFactor; i++ {
		replIdx := (index + i) % len(r.NodeHashes)
		fmt.Println(replIdx)
		nodeHash := r.NodeHashes[replIdx]
		node := r.NodeMap[nodeHash]
		nodesToReplicateTo = append(nodesToReplicateTo, node)
		fmt.Println(node.Hash)
	}

	return nodesToReplicateTo
}

// func (r *Ring) AddNode(node *Node) {
// 	r.Nodes[node.Hash] = node
// 	nodeHashes := append(r.NodeHashes, node.Hash)
// 	sortedNodeHashes := Sort(nodeHashes)
// 	r.NodeHashes = sortedNodeHashes
// }

// func (r *Ring) RemoveNode(node *Node) {
// 	index := r.Search(node.Hash)
// 	indexplus := index + 1
// 	nodeHashes := append(r.NodeHashes[:index], r.NodeHashes[index+1:])
// 	r.NodeHashes = nodeHashes
// 	delete(r.Nodes, nodeHashes[index])
// }
