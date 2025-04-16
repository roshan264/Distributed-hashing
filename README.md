# 🧩 Distributed Hashing

This is **distributed hashmap system** built across multiple servers.

## 🔧 Features

1. **Consistent Hashing**  
   Used to map keys to specific nodes, making the distribution stable and balanced across servers. 

2. **Open Addressing with Robin Hood Hashing**  
   Robin Hood hashing helps to handler collisions within O(1) lookup time.

3. **Thread Pool for HashMap Operations**  
   A worker pool manages concurrent access to the hashmap for `GET`, `PUT`, and `DELETE` operations.

4. **Flexible Value Types**  
   - **Key**: Always a string  
   - **Value**: Can be of type.

## ⚠️ Limitations (Current)

1. Nodes **cannot be added or removed** at runtime.
2. **No replication** support — value exist on only one node.

## 🌱 Planned Improvements

1. **Zookeeper Integration**  
   Node IPs/URLs will be managed through **Zookeeper** to maintain consistency.

2. **Replication Support**  
   Data will be replicated across nodes to take care of fault tolerance.

3. **Dynamic Scaling**  
   Add or remove nodes at runtime will be supported.

## 🧪 How to Test

Run the test file:

```bash
go run testing/test/test.go
