package handler

import(
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"sync"
	"distributed-hashing/server/hashmap/robinhood"
	"distributed-hashing/server/logger"
)
var LOG = logger.InitLogger("Logs/server.log")

var store sync.Map
var hm *robinhood.HashMap


func init(){
	hm = robinhood.CreateNewHashMap(0.75, 16) 
}

var pool *WorkerPool
func CreateNewWorkerPool(workerCount int, hm *robinhood.HashMap) *WorkerPool{
	pool = &WorkerPool{
		Tasks : make(chan Task),
		WorkerCount : workerCount,
		HM : hm,
	} 
	for i := 0 ; i < workerCount ; i++ {
		go pool.startWorker(i)
	}
	return pool
}

func (pool *WorkerPool) startWorker(workerId int){

	for task := range pool.Tasks {
		switch task.Operation {
			case "SET":
				err := hm.Put(task.Key, task.Value)

				if err != nil {
					task.Err <- fmt.Errorf("failed to store key %s: %v", task.Key, err)
				}else{
					task.Result <- fmt.Sprintf("Key %v and value %v are added to hashmap.", task.Key, task.Value)
				}
			case "GET":
				data, err := pool.HM.Get(task.Key)
				if err != nil {
					task.Err <- fmt.Errorf("Key %s not found : %v", task.Key, err)
				}else{
					task.Result <- data
				}
			case "DELETE":
				_, err := pool.HM.Get(task.Key)

				if err != nil{
					task.Err <- fmt.Errorf("Key %s not found : %v", task.Key, err)
				}else{
					err = pool.HM.Delete(task.Key)

					if err != nil{
						task.Err <- fmt.Errorf("Key %v could not delete due to : %v", task.Key, err)
					}else{
						task.Result <- fmt.Sprintf("Key: %v deleted from hashmap.", task.Key)
					}
				}
			default:
				task.Err <- fmt.Errorf("unknown operation: %s", task.Operation)

		}
	}
}


func (pool *WorkerPool) AddTask(task Task){
	pool.Tasks <- task 
}

func CreateHandler(port string){
	pool = CreateNewWorkerPool(20, hm)

	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/delete", handleDelete)
	address := ":" + port
	fmt.Printf("Listening on address: %v", address)
	LOG.Info("Listening on", "address", address)
	log.Fatal(http.ListenAndServe(address, nil))
}


func handleSet(w http.ResponseWriter, r *http.Request){
	var req keyValRequest
	err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Key == ""{
		LOG.Error(err, "Invalid json or missing key", "key", req.Key)
        http.Error(w, "Invalid json or missing key", http.StatusBadRequest)
        return
    }
	
	LOG.Info("Storing", "key", req.Key)
	//store.Store(req.Key, req.Value)
	//craete task for this put req
	// err = hm.Put(req.Key, req.Value)

	task := Task{
		Operation : "SET",
		Key : req.Key,
		Value : req.Value,
		Result : make(chan interface{}),
		Err : make(chan error),
	}

	LOG.Info("Adding task in pool for inserting", "key", req.Key)
	pool.AddTask(task)

	select{
	case <- task.Result:
		LOG.Info("Key", req.Key, " is stored in hashmap")
		w.WriteHeader(http.StatusCreated)
	case err := <- task.Err:
		http.Error(w, err.Error() , http.StatusInternalServerError)
		LOG.Error(err, "Error while storing ", "key", req.Key)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")

	if key == ""{
		LOG.Error(nil, "missing key", "key", key)
        http.Error(w, "Missing key", http.StatusBadRequest)
        return
    }
	LOG.Info("Fetchin", "key", key)
	task := Task{
		Operation : "GET",
		Key : key,
		Result : make(chan interface{}),
		Err : make(chan error),
	}
	LOG.Info(" for fetching", "key", key)
	pool.AddTask(task)

	select{
	case data := <- task.Result:
		val := data.([]byte)
		if !json.Valid(val) {
			http.Error(w, fmt.Sprintf("Invalid JSON stored for key: %v", key), http.StatusInternalServerError)
			return
		}
		LOG.Info("successfully got value from hashmap for ", "key", key)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(val)
	case err := <- task.Err:
		LOG.Error(err, "Error while Fetching ", "key", key)
		http.Error(w, err.Error() , http.StatusNotFound)
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")
	if key == ""{
		LOG.Error(nil, "missing key", "key", key)
        http.Error(w, "Missing key", http.StatusBadRequest)
        return
    }
	LOG.Info("Deleting", "key", key)
	task := Task{
		Operation : "DELETE",
		Key : key,
		Result : make(chan interface{}),
		Err : make(chan error),
	}
	LOG.Info("Adding task in pool for deleting", "key", key)
	pool.AddTask(task)

	select{
	case result := <- task.Result:
		LOG.Info("successfully deleted key from hashmap for ", "key", key)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, " Deleted: key %v %s", key, result)
	case err := <- task.Err:
		LOG.Error(err, "Error while Deleting ", "key", key)
		http.Error(w, err.Error() , http.StatusNotFound)
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
