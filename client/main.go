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
)

var ring *hashring.HashRing
var nodeTourlMaps = map[string]string{
	"hypervm-1":"http://localhost:9001",
	"hypervm-2":"http://localhost:9002",
	"hypervm-3":"http://localhost:9003",
}


func main(){

	ring = hashring.CreateNewHashRing()
	
	for nodeName := range nodeTourlMaps{
		ring.AddNode(nodeName)
	}

}

func setKeyValue(key string, value string){
	node := ring.Getnode(key)
	nodeUrl := nodeTourlMaps[node] + "/set"

	kv := kvRequest{Key: key, Value: value}
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

func getValue(key string) string{
	node := ring.Getnode(key)
	nodeUrl := nodeTourlMaps[node] + "/get?key=" + key
	resp, err := https.Get(url)
	if err != nil {
		return err
	}
	fmt.Printf("Response for store key:%v  is %v", key, resp)

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

