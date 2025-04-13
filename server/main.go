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
)

var store sync.Map

type keyValRequest struct {
	Key   string `json:"key"`
	Value string `json:"value",omitempty`
}

func main(){
	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/get", handleGet)

	port := "9001"
	if len(os.Args) > 1 {
		port = os.Args[1]
	} 
	address := ":" + port
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
	store.Store(req.Key, req.Value)
}

func handleGet(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")
	if value, ok := store.Load(key); ok {
		fmt.Fprintf(w, "%s", value)
	} else {
		http.NotFound(w, r)
	}

	fmt.Printf(key)
}