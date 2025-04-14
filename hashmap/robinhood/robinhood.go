package robinhood

import (
	"fmt"
	"hash/fnv"
	"encoding/json"
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
}

func CreateNewHashMap(maxLoadFactor float64, defaultCapacity int) *HashMap{
	hashMap := HashMap{
		table :make([]*entry, defaultCapacity),
		loadFactor : maxLoadFactor,
	}

	return &hashMap
}

func hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (h *HashMap) Put(key string, value interface{}) error{
	if float64(h.size + 1) / float64(len(h.table)) > h.loadFactor {
		h.resize()
	}

	valueBytes, err := json.Marshal(value)

	if err != nil{
		fmt.Printf("Error while converting value to json bytes %v", err)
		return err
	}

	newEntryAdress := &entry{Key:key, Value:valueBytes, Tombstone:false, Dist:0}

	hashVal := hash(key)
	ind := int(hashVal) %  len(h.table)

	for {
		curr := h.table[ind]

		if curr == nil || curr.Tombstone {
			h.table[ind] = newEntryAdress
			h.size++
			return nil
		}

		if curr.Key == key {
			newEntryAdress.Dist = curr.Dist 
			h.table[ind] = newEntryAdress
			return nil
		}

		if curr.Dist < newEntryAdress.Dist{
			h.table[ind], newEntryAdress = newEntryAdress, h.table[ind]
		}
		newEntryAdress.Dist++

		ind = (ind + 1 ) % len(h.table)
		
	}
}

func (h *HashMap) Get(key string) ([]byte, error){
	hashVal := hash(key)
	ind := int(hashVal) %  len(h.table) 
	
	curr := h.table[ind]
	dist := 0

	for{
		if curr == nil{
			err := fmt.Errorf("key %v not found", key)
			return nil, err
		}

		if curr.Key == key && curr.Tombstone == false{
			return curr.Value, nil
		}

		if dist > curr.Dist {
			err := fmt.Errorf("key %v not found", key)
			return nil, err
		}

		dist++
		ind = (ind + 1 ) % len(h.table)
	}
}


func (h *HashMap) Delete(key string) error {
	hashVal := hash(key)
	ind := int(hashVal) %  len(h.table) 
	
	curr := h.table[ind]
	dist := 0

	for{
		if curr == nil{
			err := fmt.Errorf("key %v not found", key)
			return err
		}

		if curr.Key == key && curr.Tombstone == false{
			curr.Tombstone = true
			h.size--
			return nil
		}

		if dist > curr.Dist {
			err := fmt.Errorf("key %v not found", key)
			return err
		}

		dist++
		ind = (ind + 1 ) % len(h.table)
	}


}

func (h *HashMap) resize() {
	newCapacity := len(h.table) * 2
	newTable := make([]*entry, newCapacity)

	oldTable := h.table 
	h.table = newTable
	h.size = 0
	for _, row := range oldTable {

		if row != nil {
			var val interface{}
			json.Unmarshal(row.Value, &val)
			h.Put(row.Key, val)
		}
	}
}

