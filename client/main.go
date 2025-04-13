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
	"github.com/roshan264/Distrubuted-Hashing/util/hashring"
	"fmt"
)

var ring *hashring.HashRing
var nodeTourlMaps = map[string]string{
	"hypervm-1":"http://localhost:9001",
	"hypervm-2":"http://localhost:9002",
	"hypervm-3":"http://localhost:9003",
}

func main(){
	ring = hashring.createNewHashRing()
	fmt.Println("ring %v", ring)
}
