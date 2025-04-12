// HashRing:{
// 	sortedHashes for nodes, nodes and hashesh mapping.
// }

// function to create hash value from key.

// AddNode function-> Add node in hashkey, also create virtual nodes and it as well.
// 					Initally assume that we are adding nodes initally only. Runtime node addition not handled.
// 				    This is critical section LOCK it. so that hashring will be consistent.

// GetNode function-> return node where key can be stored/fetched.
// 		we will have to lock this as well. But paralle READ should work when there is no WRITE going on. 
