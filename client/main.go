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
	"distributed-hashing/util/logger"
)

var ring *hashring.HashRing
var log = logger.InitLogger("/Users/StartupUser/Desktop/roshan-coding/log/client.log")

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
				log.Info("Saved: ", "key", key)
			}else{
				log.Error(err, "Could not store key")
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
				log.Info("Got it->", "key", key , "value", value)
			}else{
				log.Error(err, "Problem while fetching ", "key", key)
			}
		}

		go getKeyGoRountine (i)
		
	}

	wg.Wait()

	for i := 1 ; i < 25 ; i=i+2 {
		wg.Add(1)
		deleteKeyGoRountine := func(i int){
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			err := deleteKey(key)

			if err == nil {
				log.Info("Deleted->", "key", key )
			}else{
				log.Error(err, "Problem while deleting ", "key", key)
			}
		}

		go deleteKeyGoRountine (i)
		
	}

	wg.Wait()

	for i := 1 ; i < 21 ; i++ {
		wg.Add(1)
		getKeyGoRountine := func(i int){
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value, err := getValue(key)

			if err == nil {
				log.Info("Got it->", "key", key , "value", value)
			}else{
				log.Error(err, "Problem while fetching ", "key", key)
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
		log.Error(err, "Failed to convert keyValRequest to json")
		return err
	}
	log.Info("calling: ", nodeUrl)
	resp, err := http.Post(nodeUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Error(err, "Failed while calling url", nodeUrl)
		return err
	}
	log.Info("Response for store key: ", key, "value: ", value, "resp: ", resp)
	defer resp.Body.Close()
	return nil
}

func getValue(key string) (string, error){
	node := ring.GetNode(key)
	if node == ""{
		err := fmt.Errorf("Unable to get node for key %s", key)
		log.Error(err, "Failed in getting node for key")
		return "", err
	}
	nodeUrl := nodeTourlMaps[node] + "/get?key=" + key
	log.Info("Calling ", "url",  nodeUrl)
	resp, err := http.Get(nodeUrl)
	if err != nil {
		log.Error(err, "Error while calling", "url",nodeUrl )
		return "", err
	}
	// fmt.Printf("Response for store key:%v  is %v", key, resp)

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound{
		err := fmt.Errorf("key not found %s", key)
		log.Error(err, "Failed to get key")
		return "", err
	}else if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Info("Resonse for key", key , "Value: ", string(body))
		return string(body), nil
	}else{
		err := fmt.Errorf("Unexecpcted status code: %v", resp.StatusCode)
		return "", err
	}
}

func deleteKey(key string) error {
	node := ring.GetNode(key)
	if node == ""{
		err := fmt.Errorf("Unable to get node for key %s", key)
		log.Error(err, "Failed in getting node for key")
		return err
	}

	nodeUrl := nodeTourlMaps[node] + "/delete?key=" + key 
	log.Info("Calling ", "url",  nodeUrl)
	
	req, err := http.NewRequest(http.MethodDelete, nodeUrl, nil)
	if err != nil {
		log.Error(err, "Failed to create DELETE request for", nodeUrl)
		return err
	}
	client := &http.Client{} 
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err, "Error while calling delete req for", "key", key)
		return err
	}

	log.Info("Response ", "status", resp.Status, "key", key)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err, "Error while deleting ", "key", key)
		return err
	}
	log.Info("Response ", "Body", string(respBody))

	return nil
}

