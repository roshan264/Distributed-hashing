package test 

import(
	"fmt"
	"strconv"
	"sync"
	"distributed-hashing/client/methods"
	"encoding/json"
)

type User struct {
	Name string
	Age  int
}


func DecodeValue[T any](data []byte) (T, error) {
	var t T
	err := json.Unmarshal(data, &t)
	return t, err
}

func UnitTesting(){

	var wg sync.WaitGroup
	fmt.Printf("\n********Adding Keys to HashMap*****\n")
	for i := 0 ; i < 8 ; i++{
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			name := "Roshan" + strconv.Itoa(i)
			key := "user" + strconv.Itoa(i)
			age := i 
			user := User{Name: name, Age: age}
			err := methods.SetKeyValue(key, user)
			fmt.Printf("Added key: %v value: %+v\n", key, user)
			if err != nil {
				fmt.Printf("Error while setting :key: %v: %v \n", key, err.Error())
			}
		}(i)
		//go putKeyGoroutine (i)
	}

	wg.Wait()
	fmt.Printf("\n******************************\n\n")
	fmt.Printf("********Fetching Keys from hashmap*****\n")
	for i := 0 ; i < 8 ; i++{
		wg.Add(1)
		getKeyGoRountine := func(i int){
			defer wg.Done()
			key := "user" + strconv.Itoa(i)
			data, err := methods.GetValue(key)
			if err != nil {
				fmt.Printf("key %v not found", key)
				return
			}
			fmt.Printf("Fetched key: %v value: %+v\n", key, string(data))
			// val, _ := DecodeValue[User](data)

			// fmt.Printf("fetching(After decoding to value format) for key: %v value: %v \n", name, val)
			// fmt.Printf("Print: name %v", val.Name)
		}
		go getKeyGoRountine(i)
	}

	wg.Wait()

	fmt.Printf("\n******************************\n\n")
	fmt.Printf("********Deleting few Keys from hashmap*****\n")

	for i := 0 ; i < 12 ; i = i + 3{
		wg.Add(1)
		deleteKeyGoRountine := func(i int){
			defer wg.Done()
			key := "user" + strconv.Itoa(i)
			err := methods.DeleteKey(key)
			if err != nil {
				fmt.Printf("Error while deleting %v : %v \n", key, err)
				return
			}

			fmt.Printf("Deleted key: %v \n", key)
			
		}
		go deleteKeyGoRountine(i)
	}

	wg.Wait()
	fmt.Printf("\n******************************\n\n")
	fmt.Printf("********Again fetching Keys from hashmap*****\n")

	for i := 0 ; i < 10; i++{
		wg.Add(1)
		getKeyGoRountine := func(i int){
			defer wg.Done()
			key := "user" + strconv.Itoa(i)
			data, err := methods.GetValue(key)
			if err != nil {
				fmt.Printf("key %v not found %v \n", key, err)
				return
			}
			fmt.Printf("Fetched: key: %v value: %+v\n", key, string(data))
			// val, _ := DecodeValue[User](data)

			// fmt.Printf("fetching(After decoding to value format) for key: %v value: %v \n", name, val)
			// fmt.Printf("Print: name %v", val.Name)
		}
		go getKeyGoRountine(i)
	}

	wg.Wait()

	fmt.Printf("\n******************************\n\n")
	fmt.Printf("********Overwrite values of some keys*****\n")

	for i := 0 ; i < 10 ; i = i + 2{
		wg.Add(1)
		overwriteKeyGoRountine := func(i int){
			defer wg.Done()
			name := "Roshan-overwrite" + strconv.Itoa(i)
			key := "user" + strconv.Itoa(i)
			age := i 
			user := User{Name: name, Age: age}
			err := methods.SetKeyValue(key, user)
			fmt.Printf("Overwrite key: %v value: %+v\n", key, user)
			if err != nil {
				fmt.Printf("Error while setting :key: %v: %v \n", key, err.Error())
			}
			
		}
		go overwriteKeyGoRountine(i)
	}
	wg.Wait()

	fmt.Printf("\n******************************\n\n")
	fmt.Printf("********Fetch Overwritten values *****\n")

	for i := 0 ; i < 10  ; i = i + 2{
		wg.Add(1)
		overwriteKeyGoRountine := func(i int){
			defer wg.Done()
			key := "user" + strconv.Itoa(i)
			data, err := methods.GetValue(key)
			if err != nil {
				fmt.Printf("key %v not found\n", key)
				return
			}
			fmt.Printf("Fetched: key: %v value: %+v\n", key, string(data))

			
		}
		go overwriteKeyGoRountine(i)
	}
	wg.Wait()
	
	fmt.Println("Test completed successfully!")
}

