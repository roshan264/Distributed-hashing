package main

import (
	"fmt"
	"log"
	"strconv"

	"distributed-hashing/hashmap/robinhood"
)


type User struct {
	Name string
	Age  int
}

func main() {
	// Create a new HashMap
	hm := robinhood.CreateNewHashMap(0.75, 16)

	// Test data for Put and Get operations
	for i := 0 ; i < 18 ; i++{

		name := "Roshan" + strconv.Itoa(i)
		age := i 
		user := User{Name: name, Age: age}
		err := hm.Put(name, user)
		fmt.Printf("post: key: %v value: %+v\n", name, age)
		if err != nil {
			log.Fatalf("Put failed: %v", err)
		}
	}

	for i := 0 ; i < 18 ; i++{

		name := "Roshan" + strconv.Itoa(i)

		data, err := hm.Get(name)
		if err != nil {
			fmt.Printf("Error while fetchi %v : %v", name, err)
		}
		fmt.Printf("without marshal: key: %v value: %+v\n", name, string(data))
	}



	data, err := hm.Get("user1")
	if err != nil {
		log.Fatalf("Get failed: %v", err)
	}
	fmt.Printf("without marshal: %+v\n", string(data))

	fmt.Println("Test completed successfully!")
}


