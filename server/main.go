// For first cut will use go's inbuilt map.
//	create two handlers get and set.

// Set-> set key value in map.-> we will have to add LOCK here so that no one can read/write value from map.

// get-> return value from map.-> multiple reader can be allowed. But when write is going on we cant read.

package main 

import(
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"sync"
	"distributed-hashing/hashmap/robinhood"
)

var store sync.Map
var hm *robinhood.HashMap
type keyValRequest struct {
	Key   string `json:"key"`
	Value interface{} `json:"value",omitempty`
}

func main(){
	fmt.Printf("roshan")
	hm = robinhood.CreateNewHashMap(0.75, 16) 
	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/delete", handleDelete)
	port := "9001"
	if len(os.Args) > 1 {
		port = os.Args[1]
	} 
	address := ":" + port
	fmt.Printf("Listening on address: %v", address)

	log.Fatal(http.ListenAndServe(address, nil))
}

func handleSet(w http.ResponseWriter, r *http.Request){
	var req keyValRequest
	err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	fmt.Printf("key:%v value:%v", req.Key, req.Value)
	//store.Store(req.Key, req.Value)
	err = hm.Put(req.Key, req.Value)
	if err != nil{
		http.Error(w, fmt.Sprintf("Failed to store key: %v Error: %v", req.Key, err), http.StatusInternalServerError)	
	}

}

func handleGet(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")
	// if value, ok := store.Load(key); ok {
	// 	fmt.Fprintf(w, "%s", value)
	// } else {
	// 	http.NotFound(w, r)
	// }
	// if err != nil {
	// 	fmt.Printf("Error while fetching %v : %v", key, err)
	// 	http.Error(w, fmt.Sprintf("Error while fetching key: %v Error: %v", key, err), http.StatusInternalServerError)	
	// }
	// w.Header().Set("Content-Type", "application/json")
	// err = json.NewEncoder(w).Encode(data)
	// if err != nil{
	// 	http.Error(w, fmt.Sprintf("Error correcting json: %v Error: %v", key, err), http.StatusInternalServerError)	
	// }
	// fmt.Fprintf(w, "%v", data)


	data, err := hm.Get(key)
	if err != nil {
		fmt.Printf("Error while fetching %v : %v\n", key, err)
		http.Error(w, fmt.Sprintf("Error while fetching key: %v Error: %v", key, err), http.StatusInternalServerError)
		return
	}

	// Validate it's proper JSON
	if !json.Valid(data) {
		http.Error(w, fmt.Sprintf("Invalid JSON stored for key: %v", key), http.StatusInternalServerError)
		
	}else{
		fmt.Printf("Valid json %v", data)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// _, ok := data.([]byte)
	// if !ok {
	// 	fmt.Printf("value is not in []byte format")
	// }else{
	// 	fmt.Printf("value is in []byte format")
	// }
	w.Write(data) // Send the raw JSON as-is
}

func handleDelete(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")

	_, err := hm.Get(key)
	if err != nil {
		fmt.Printf("Error while fetching %v : %v\n", key, err)
		http.Error(w, fmt.Sprintf("Error while fetching key: %v Error: %v", key, err), http.StatusInternalServerError)
		return
	}
	
	err = hm.Delete(key)

	if err != nil {
		fmt.Printf("Error while fetching %v : %v\n", key, err)
		http.Error(w, fmt.Sprintf("Error while fetching key: %v Error: %v", key, err), http.StatusInternalServerError)
		return
	}else{
		fmt.Fprintf(w, "Deleted: %s", key)
	}

	// _, ok := store.Load(key)
	// if ok {
	// 	store.Delete(key)
	// 	fmt.Fprintf(w, "Deleted: %s", key)
	// } else {
	// 	fmt.Printf("Key: %v node not present", key)
	// 	http.NotFound(w, r)
	// }

}