package handler

import (
	"sync"

	"distributed-hashing/server/hashmap/robinhood"
)
type keyValRequest struct {
	Key   string `json:"key"`
	Value interface{} `json:"value",omitempty`
}

type Task struct{
	ID string 
	Operation string
	Key string 
	Value  interface{}
	Result chan interface {}
	Err chan error
}

type WorkerPool struct{
	Tasks chan Task
	WorkerCount int 
	HM *robinhood.HashMap
	Wg sync.WaitGroup
}
