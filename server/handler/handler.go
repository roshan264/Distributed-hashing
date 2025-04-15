package handler

import(
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"sync"
	"distributed-hashing/server/hashmap/robinhood"
)

var store sync.Map
var hm *robinhood.HashMap
type keyValRequest struct {
	Key   string `json:"key"`
	Value interface{} `json:"value",omitempty`
}

func init(){
	hm = robinhood.CreateNewHashMap(0.75, 16) 
}

func CreateHandler(port string){

	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/delete", handleDelete)
	address := ":" + port
	fmt.Printf("Listening on address: %v", address)

	log.Fatal(http.ListenAndServe(address, nil))
}


func handleSet(w http.ResponseWriter, r *http.Request){
	var req keyValRequest
	err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Key == ""{
        http.Error(w, "Invalid json or missing key", http.StatusBadRequest)
        return
    }
	fmt.Printf("key:%v value:%v", req.Key, req.Value)
	//store.Store(req.Key, req.Value)
	err = hm.Put(req.Key, req.Value)
	if err != nil{
		http.Error(w, fmt.Sprintf("Failed to store key: %v Error: %v", req.Key, err), http.StatusInternalServerError)	
	}
	w.WriteHeader(http.StatusCreated)
}

func handleGet(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")

	if key == ""{
        http.Error(w, "Missing key", http.StatusBadRequest)
        return
    }

	data, err := hm.Get(key)
	if err != nil {
		// fmt.Printf("Error while fetching %v : %v\n", key, err)
		http.Error(w, fmt.Sprintf("Error while fetching key: %v Error: %v", key, err), http.StatusInternalServerError)
		return
	}

	if !json.Valid(data) {
		http.Error(w, fmt.Sprintf("Invalid JSON stored for key: %v", key), http.StatusInternalServerError)	
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data) 
}

func handleDelete(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")
	if key == ""{
        http.Error(w, "Missing key", http.StatusBadRequest)
        return
    }
	_, err := hm.Get(key)

	if err != nil {
		//fmt.Printf("Error while fetching %v : %v\n", key, err)
		http.Error(w, fmt.Sprintf("Error while fetching key: %v Error: %v", key, err), http.StatusInternalServerError)
		return
	}
	
	err = hm.Delete(key)

	if err != nil {
		fmt.Printf("Error while fetching %v : %v\n", key, err)
		http.Error(w, fmt.Sprintf("Error while Deleting key: %v Error: %v", key, err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted: %s", key)


	// _, ok := store.Load(key)
	// if ok {
	// 	store.Delete(key)
	// 	fmt.Fprintf(w, "Deleted: %s", key)
	// } else {
	// 	fmt.Printf("Key: %v node not present", key)
	// 	http.NotFound(w, r)
	// }

}