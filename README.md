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

1. Clone the repo. And run make command like below in repo directory. 

```bash
make run-all

# Set data
curl -X POST "http://localhost:9004/set?key=user123" \
     -H "Content-Type: application/json" \
     -d '{
           "data": {
               "name": "roshan",
               "age": 25,
               "tags": ["go", "dev"],
               "meta": {
                   "active": true,
                   "lastLogin": "2024-12-01T12:34:56Z"
               }
           }
         }'

# Get data
curl -X GET "http://localhost:9004/get?key=user123"

# Delete data
curl -X DELETE "http://localhost:9004/get?key=user123"


