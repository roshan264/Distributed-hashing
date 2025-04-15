package methods

import (
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
	"io"
	// "reflect"

	"distributed-hashing/util/hashring"
	"distributed-hashing/util/logger"
)

var log = logger.InitLogger("/Users/StartupUser/Desktop/roshan-coding/log/client.log")

var NodeTourlMaps = map[string]string{
	"hypervm-1":"http://localhost:9001",
	"hypervm-2":"http://localhost:9002",
	"hypervm-3":"http://localhost:9003",
}

type keyValRequest struct {
	Key   string `json:"key"`
	Value interface{} `json:"value",omitempty`
}

func SetKeyValue(key string, value interface{}, ring *hashring.HashRing) error {
	node := ring.GetNode(key)
	nodeUrl := NodeTourlMaps[node] + "/set"

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

func GetValue(key string, ring *hashring.HashRing) ([]byte, error){
	node := ring.GetNode(key)
	if node == ""{
		err := fmt.Errorf("Unable to get node for key %s", key)
		log.Error(err, "Failed in getting node for key")
		return nil, err
	}
	nodeUrl := NodeTourlMaps[node] + "/get?key=" + key
	log.Info("Calling ", "url",  nodeUrl)
	resp, err := http.Get(nodeUrl)
	if err != nil {
		log.Error(err, "Error while calling", "url",nodeUrl )
		return nil, err
	}
	// fmt.Printf("Response for store key:%v  is %v", key, resp)

	defer resp.Body.Close()

	// if resp.StatusCode == http.StatusNotFound{
	// 	err := fmt.Errorf("key not found %s", key)
	// 	log.Error(err, "Failed to get key")
	// 	return nil, err
	// }else if resp.StatusCode == http.StatusOK {
	// 	var value interface{}
	// 	body, _ := io.ReadAll(resp.Body)
	// 	log.Info("body", body, "roshan")
	// 	fmt.Printf("body: %v",body)
	// 	//err := json.Unmarshal(body, &value)
	// 	bodyStr := string(body)
	// 	fmt.Printf(bodyStr)
	// 	err = json.Unmarshal([]byte(body), &value)
	// 	if err != nil {
	// 		log.Error(err, "Failed to decode JSON:")
	// 		return nil, err
	// 	}
	// 	log.Info("Resonse for key", key , "Value: ", value)
	// 	return value, nil
	// }else{
	// 	err := fmt.Errorf("Unexecpcted status code: %v", resp.StatusCode)
	// 	return "", err
	// }


	if resp.StatusCode == http.StatusNotFound {
		err := fmt.Errorf("key not found %s", key)
		log.Error(err, "Failed to get key")
		return nil, err
	} else if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "Failed to read response body")
			return nil, err
		}
		// if reflect.TypeOf(body) == reflect.TypeOf([]byte{}) {
		// 	fmt.Printf("It's []byte: %v", key)
		// }
		log.Info("Raw JSON body received", "key", key, "body", string(body))
		return body, nil // Return as []byte
	} else {
		err := fmt.Errorf("Unexpected status code: %v", resp.StatusCode)
		return nil, err
	}
}

func DeleteKey(key string, ring *hashring.HashRing) error {
	node := ring.GetNode(key)
	if node == ""{
		err := fmt.Errorf("Unable to get node for key %s", key)
		log.Error(err, "Failed in getting node for key")
		return err
	}

	nodeUrl := NodeTourlMaps[node] + "/delete?key=" + key 
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

