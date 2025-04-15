package test 

import(
	"fmt"
	"strconv"
	"sync"
	"log"
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

	for i := 0 ; i < 58 ; i++{
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			name := "Roshan" + strconv.Itoa(i)
			age := i 
			user := User{Name: name, Age: age}
			err := methods.SetKeyValue(name, user)
			fmt.Printf("post: key: %v value: %+v\n", name, age)
			if err != nil {
				log.Fatalf("Put failed: %v", err)
			}
		}(i)
		//go putKeyGoroutine (i)
	}

	wg.Wait()

	for i := 0 ; i < 58 ; i++{
		wg.Add(1)
		getKeyGoRountine := func(i int){
			defer wg.Done()
			name := "Roshan" + strconv.Itoa(i)

			data, err := methods.GetValue(name)
			if err != nil {
				fmt.Printf("Error while fetchi %v : %v", name, err)
			}
			// fmt.Printf("without marshal: key: %v value: %+v\n", name, string(data))
			val, _ := DecodeValue[User](data)

			fmt.Printf("fetched : key: %v value: %v\n", name, val)
			// fmt.Printf("Print: name %v", val.Name)
		}
		go getKeyGoRountine(i)
	}

	wg.Wait()


	for i := 0 ; i < 70 ; i = i + 4{
		wg.Add(1)
		deleteKeyGoRountine := func(i int){
			defer wg.Done()
			name := "Roshan" + strconv.Itoa(i)

			err := methods.DeleteKey(name)
			if err != nil {
				fmt.Printf("Error while deleting %v : %v", name, err)
			}

			fmt.Printf("delted key key: %v \n", name)
			
		}
		go deleteKeyGoRountine(i)
	}

	wg.Wait()


	for i := 0 ; i < 58 ; i++{
		wg.Add(1)
		getKeyGoRountine := func(i int){
			defer wg.Done()
			name := "Roshan" + strconv.Itoa(i)

			data, err := methods.GetValue(name)
			if err != nil {
				fmt.Printf("Error while fetching %v : %v", name, err)
			}
			// fmt.Printf("without marshal: key: %v value: %+v\n", name, string(data))
			val, _ := DecodeValue[User](data)

			fmt.Printf("Fetched key: %v value: %v\n", name, val)
			// fmt.Printf("Print: name %v", val.Name)
		}
		go getKeyGoRountine(i)
	}

	wg.Wait()

	fmt.Println("Test completed successfully!")
}

