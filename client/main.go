/*store server -> url maps.

add one by one node from server in hashring.
give few set key-value calls.
also fetch those values.

set-key value function: 
	convert key into hash.
	call hashring get node function.
	send key on given node.
get key value function:
	same as above.
delete key same as above.*/


package main

import(
	"distributed-hashing/util/logger"
	"distributed-hashing/client/methods"
	"distributed-hashing/testing/test"
)

var log = logger.InitLogger("/Users/StartupUser/Desktop/roshan-coding/log/client.log")


func main(){
	methods.Setup()
	test.UnitTesting()
}

