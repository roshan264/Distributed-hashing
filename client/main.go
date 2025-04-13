/*store server -> url maps.

add one by one node from server in hashring.
give few set key-value calls.
also fetch those values.

set-key value function: 
	convert key into hash.
	call hashring get node function.
	send key on given node.
get key value function:
	same as above.
delete key same as above.*/


package main

import(
	"distributed-hashing/util/hashring"
	"fmt"
	"strconv"
	"sync"
	"encoding/json"
	"net/http"
	"bytes"
	"io"
)

var ring *hashring.HashRing
var nodeTourlMaps = map[string]string{
	"hypervm-1":"http://localhost:9001",
	"hypervm-2":"http://localhost:9002",
	"hypervm-3":"http://localhost:9003",
}

type keyValRequest struct {
	Key   string `json:"key"`
	Value string `json:"value",omitempty`
}

func main(){

	ring = hashring.CreateNewHashRing()
	
	for nodeName := range nodeTourlMaps{
		ring.AddNode(nodeName)
	}

	var wg sync.WaitGroup

	for i := 1 ; i < 21 ; i++ {
		wg.Add(1)
		setKeyGoRountine := func(i int){
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value := "value" + strconv.Itoa(i)
			err := setKeyValue(key, value)

			if err == nil {
				fmt.Printf("Saved: key:%v", key)
			}
		}

		go setKeyGoRountine (i)
		
	}

	wg.Wait()

	for i := 1 ; i < 21 ; i++ {
		wg.Add(1)
		getKeyGoRountine := func(i int){
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value, err := getValue(key)

			if err == nil {
				fmt.Printf("Got it-> key %v : value %v", key, value)
			}else{
				fmt.Printf("Error while fetchin value for %v : %v", key, err)
			}
		}

		go getKeyGoRountine (i)
		
	}

	wg.Wait()

}

func setKeyValue(key string, value string) error {
	node := ring.GetNode(key)
	nodeUrl := nodeTourlMaps[node] + "/set"

	kv := keyValRequest{Key: key, Value: value}
	data, err := json.Marshal(kv)
	if err != nil {
		return err
	}

	resp, err := http.Post(nodeUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	fmt.Printf("Response for store key:%v and value:%v is %v", key, value, resp)
	defer resp.Body.Close()
	return nil
}

func getValue(key string) (string, error){
	node := ring.GetNode(key)
	nodeUrl := nodeTourlMaps[node] + "/get?key=" + key
	resp, err := http.Get(nodeUrl)
	if err != nil {
		return "", err
	}
	fmt.Printf("Response for store key:%v  is %v", key, resp)

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("getValue Body: %v", body)
	return string(body), nil
}

