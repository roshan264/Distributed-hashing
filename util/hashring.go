// HashRing:{
// 	sortedHashes for nodes, nodes and hashesh mapping.
// }

// function to create hash value from key.

// AddNode function-> Add node in hashkey, also create virtual nodes and it as well.
// 					Initally assume that we are adding nodes initally only. Runtime node addition not handled.
// 				    This is critical section LOCK it. so that hashring will be consistent.

// GetNode function-> return node where key can be stored/fetched.
// 		we will have to lock this as well. But paralle READ should work when there is no WRITE going on. 


package hashring

import(
	"sort"
	"sync"
	"strconv"
	"crypto/sha256"
)
int virtualNodes = 100

type HashRing struct{
	nodes map[uint64]string 
	sortedHashes []uint64
	mutex sync.RWMutex
}

func createNewHashRing() *HashRing{
	hashRing := HashRing{
		nodes: make(map[uint64]string)
		sortedHashes: make([]uint64, 0)
	}
	return &hashRing
}

func convertKeyToHash(key string){
	sum := sha256.Sum256([]byte(key))
	return uint64(sum[0])<<56 | uint64(sum[1])<<48 | uint64(sum[2])<<40 | uint64(sum[3])<<32 |
		uint64(sum[4])<<24 | uint64(sum[5])<<16 | uint64(sum[6])<<8 | uint64(sum[7])
}


func (h *HashRing) addNode(nodeName string){

	h.mutex.Lock()
	defer h.mutex.UnLock()

	for i := 0 ; i <= virtualNodes ; i++{
		virtualNodeId := nodeName + "_" + i
		hashKeyOfVirtualNode := hashKey(virtualNodeId)
		h.nodes[hashKeyOfVirtualNode] = nodeName
		h.sortedHashes = append(h.sortedHashes, hashKeyOfVirtualNode)
	}

	sort.Slice(h.sortedHashes, func(i, j int) bool{
		return h.sortedHashes[i] < h.sortedHashes[j]
	})

}

func FindTargetedNodeHash(sortedHashes uint64[], hashKey uint64){
	length := len(sortedHashes)
	if length == 0{
		return ""
	}
	if length == 1{
		return sortedHashes[0]
	}

	
	left := 0
	right := length - 1

	for left <= right {
		mid := left + (right - left) / 2

		if arr[mid] <= key {
			left = mid + 1
		}else{
			right = mid - 1
		}
	}

	if left < length{
		return arr[left]
	}

	return arr[0]
	
}
func (h *HashRing) getNode(key string) string{
	h.mutex.Rlock()
	defer h.mutex.RUnlock()

	if len(h.sortedHashes) == 0{
		return ""
	}

	hashKey = convertKeyToHash(key)

	targetedHash = FindTargetedNodeHash(h.sortedHashes, hashKey)

	if targetedHash != ""{
		return h.nodes[targetedHash]
	}

	return ""
}

