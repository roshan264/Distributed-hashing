package robinhood

import(
	"sync"
)
type entry struct {
	Key string
	Value []byte
	Tombstone bool 
	Dist int
} 


type HashMap struct{
	table    []*entry
	size int 
	loadFactor float64
	mutex sync.RWMutex
}