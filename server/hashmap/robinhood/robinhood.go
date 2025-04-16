package robinhood

import (
	"fmt"
	"hash/fnv"
	"encoding/json"

	"distributed-hashing/server/logger"
)
var LOG = logger.InitLogger("Logs/hashmap.log")

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

func (h *HashMap) putInternal(key string, value interface{}, convertVal bool) error{
	valueBytes, err := json.Marshal(value)

	if err != nil{
		LOG.Error(err, "Error while converting value to json bytes")
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

	LOG.Error(nil, "Execution should never reach here.")
	return nil

}
func (h *HashMap) Put(key string, value interface{}) error{
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if float64(h.size + 1) / float64(len(h.table)) > h.loadFactor {
		h.resize()
	}

	// h.PrintMap()
	return h.putInternal(key, value, true)
}

func (h *HashMap) Get(key string) ([]byte, error){
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	hashVal := hash(key)
	ind := int(hashVal) %  len(h.table) 
	
	
	dist := 0

	for{
		curr := h.table[ind]
		if curr == nil{
			err := fmt.Errorf("key %v not found", key)
			LOG.Info(err.Error())
			return nil, err
		}

		if curr.Key == key && curr.Tombstone == false{
			return curr.Value, nil
		}

		if dist > curr.Dist {
			err := fmt.Errorf("key %v not found", key)
			LOG.Info(err.Error())
			return nil, err
		}

		dist++
		ind = (ind + 1 ) % len(h.table)
	}
}


func (h *HashMap) Delete(key string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	hashVal := hash(key)
	ind := int(hashVal) %  len(h.table) 
	
	
	dist := 0

	for{
		curr := h.table[ind]
		if curr == nil{
			err := fmt.Errorf("key %v not found", key)
			LOG.Error(err, "")
			return err
		}

		if curr.Key == key && curr.Tombstone == false{
			curr.Tombstone = true
			h.size--
			return nil
		}

		if dist > curr.Dist {
			err := fmt.Errorf("key %v not found", key)
			LOG.Error(err, "")
			return err
		}

		dist++
		ind = (ind + 1 ) % len(h.table)
	}

	return nil

}

func (h *HashMap) resize() {
	//This function also needs lock. But we are calling this from Put only. And put already has lock. So not locking here.
	newCapacity := len(h.table) * 2
	newTable := make([]*entry, newCapacity)

	oldTable := h.table 
	h.table = newTable
	h.size = 0
	for _, row := range oldTable {

		if row != nil && !row.Tombstone {
			var val interface{}
			json.Unmarshal(row.Value, &val)
			h.putInternal(row.Key, val , false)
		}
	}
}

func (h *HashMap) PrintMap(){
	for _, row := range h.table{
		if row != nil {
			LOG.Info("key" , row.Key, "Distanve", row.Dist, "Tombstone", row.Tombstone)
		}
		
	}
}