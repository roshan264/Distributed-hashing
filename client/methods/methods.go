package methods

import (
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
	"io"
	// "reflect"

	"distributed-hashing/client/util/hashring"
	"distributed-hashing/client/util/logger"
)

var LOG = logger.InitLogger("Logs/client.log")

var NodeTourlMaps = map[string]string{
	"hypervm-1":"http://localhost:9001",
	"hypervm-2":"http://localhost:9002",
	"hypervm-3":"http://localhost:9003",
}

type keyValRequest struct {
	Key   string `json:"key"`
	Value interface{} `json:"value",omitempty`
}
var ring *hashring.HashRing

func Setup() {

	ring = hashring.CreateNewHashRing()	
	for nodeName := range NodeTourlMaps{
		ring.AddNode(nodeName)
	}
}

func SetKeyValue(key string, value interface{}) error {
	node := ring.GetNode(key)
	nodeUrl := NodeTourlMaps[node] + "/set"

	kv := keyValRequest{Key: key, Value: value}
	data, err := json.Marshal(kv)
	if err != nil {
		LOG.Error(err, "Failed to convert keyValRequest to json")
		return err
	}
	LOG.Info("calling: ", nodeUrl)
	resp, err := http.Post(nodeUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		LOG.Error(err, "Failed while calling url", nodeUrl)
		return err
	}
	LOG.Info("Response for store key: ", key, "resp: ", resp)
	defer resp.Body.Close()
	return nil
}

func GetValue(key string) ([]byte, error){
	node := ring.GetNode(key)
	if node == ""{
		err := fmt.Errorf("Unable to get node for key %s", key)
		LOG.Error(err, "Failed in getting node for key")
		return nil, err
	}
	nodeUrl := NodeTourlMaps[node] + "/get?key=" + key
	LOG.Info("Calling ", "url",  nodeUrl)
	resp, err := http.Get(nodeUrl)
	if err != nil {
		LOG.Error(err, "Error while calling", "url",nodeUrl )
		return nil, err
	}
	// fmt.Printf("Response for store key:%v  is %v", key, resp)

	defer resp.Body.Close()

	// if resp.StatusCode == http.StatusNotFound{
	// 	err := fmt.Errorf("key not found %s", key)
	// 	LOG.Error(err, "Failed to get key")
	// 	return nil, err
	// }else if resp.StatusCode == http.StatusOK {
	// 	var value interface{}
	// 	body, _ := io.ReadAll(resp.Body)
	// 	LOG.Info("body", body, "roshan")
	// 	fmt.Printf("body: %v",body)
	// 	//err := json.Unmarshal(body, &value)
	// 	bodyStr := string(body)
	// 	fmt.Printf(bodyStr)
	// 	err = json.Unmarshal([]byte(body), &value)
	// 	if err != nil {
	// 		LOG.Error(err, "Failed to decode JSON:")
	// 		return nil, err
	// 	}
	// 	LOG.Info("Resonse for key", key , "Value: ", value)
	// 	return value, nil
	// }else{
	// 	err := fmt.Errorf("Unexecpcted status code: %v", resp.StatusCode)
	// 	return "", err
	// }


	if resp.StatusCode == http.StatusNotFound {
		err := fmt.Errorf("key not found %s", key)
		LOG.Error(err, "Failed to get key")
		return nil, err
	} else if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			LOG.Error(err, "Failed to read response body")
			return nil, err
		}
		// if reflect.TypeOf(body) == reflect.TypeOf([]byte{}) {
		// 	fmt.Printf("It's []byte: %v", key)
		// }
		LOG.Info("Raw JSON body received", "key", key, "body", string(body))
		return body, nil // Return as []byte
	} else {
		err := fmt.Errorf("Key %v not found, Unexpected status code: %v", key, resp.StatusCode)
		return nil, err
	}
}

func DeleteKey(key string) error {
	node := ring.GetNode(key)
	if node == ""{
		err := fmt.Errorf("Unable to get node for key %s", key)
		LOG.Error(err, "Failed in getting node for key")
		return err
	}

	nodeUrl := NodeTourlMaps[node] + "/delete?key=" + key 
	LOG.Info("Calling ", "url",  nodeUrl)
	
	req, err := http.NewRequest(http.MethodDelete, nodeUrl, nil)
	if err != nil {
		LOG.Error(err, "Failed to create DELETE request for", nodeUrl)
		return err
	}
	client := &http.Client{} 
	resp, err := client.Do(req)
	if err != nil {
		LOG.Error(err, "Error while calling delete req for", "key", key)
		return err
	}

	LOG.Info("Response ", "status", resp.Status, "key", key)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		LOG.Error(err, "Error while deleting ", "key", key)
		return err
	}
	LOG.Info("Response ", "Body", string(respBody))

	return nil
}

