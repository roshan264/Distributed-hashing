package handler

import(
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"distributed-hashing/client/util/logger"
	"distributed-hashing/client/methods"
)
var LOG = logger.InitLogger("Logs/client.log")


var setCmd = `curl -X POST "http://localhost:9004/set?key=user123" \
     -H "Content-Type: application/json" \
     -d '{
           "data": {
               "name": "roshan",
               "age": 25,
               "tags": ["go", "dev"],
               "meta": {
                   "active": true,
                   "lastLogin": "2024-12-01T12:34:56Z"
               }
           }
         }'`
var getCmd = `curl -X GET "http://localhost:9004/get?key=user123"`
var deleteCmd = `curl -X DELETE "http://localhost:9004/delete?key=user123"`

func CreateHandler(port string){	
	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/delete", handleDelete)
	address := ":" + port
	fmt.Printf("\nListening on address: %v\n", address)

	fmt.Printf("Format to store key-value pair command\n %v\n \n", setCmd)
	fmt.Printf("Format to get key-value pair command\n %v\n \n", getCmd)
	fmt.Printf("Format to delete key-value pair command\n %v\n \n", deleteCmd)

	LOG.Info("Listening on", "address", address)
	log.Fatal(http.ListenAndServe(address, nil))
}


func handleSet(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")

	if key == ""{
		LOG.Error(nil, "missing key", "key", key)
        http.Error(w, "Missing key", http.StatusBadRequest)
        return
    }
	
	var value interface{}
	err := json.NewDecoder(r.Body).Decode(&value)
	if err != nil{
		LOG.Error(err, "Error while decoding value", "value", value)
        http.Error(w, "Error while decoding value", http.StatusBadRequest)
        return
    }

	err = methods.SetKeyValue(key, value)
	LOG.Info("Added key: %v value: %v\n", key, value)
	if err != nil {
		fmt.Printf("Error while setting :key: %v: %v \n", key, err.Error())
	}

	LOG.Info("successfully stored key in hashmap Key: ", "key", key)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, " Stored: key %v ", key)
}

func handleGet(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")

	if key == ""{
		LOG.Error(nil, "missing key", "key", key)
        http.Error(w, "Missing key", http.StatusBadRequest)
        return
    }
	
	data, err := methods.GetValue(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	LOG.Info("successfully got value from hashmap for ", "key", key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleDelete(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")

	if key == ""{
		LOG.Error(nil, "missing key", "key", key)
        http.Error(w, "Missing key", http.StatusBadRequest)
        return
    }

	err := methods.DeleteKey(key)
	if err != nil {
		fmt.Printf("Error while deleting %v : %v \n", key, err)
		http.Error(w, err.Error() , http.StatusNotFound)
		return
	}
	LOG.Info("successfully deleted key from hashmap Key: ", "key", key)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, " Deleted: key %v", key)
}
