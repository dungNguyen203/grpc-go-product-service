- Server: Provides gRPC methods with metadata-based authentication, manages Product data using GORM and SQLite.

- gRPC-Gateway: Exposes REST endpoints, mapping directly to gRPC methods for seamless HTTP/2 and REST support.

- Client: Using Auth Interceptor, invokes CreateProduct and streams ListProducts.

- Gin REST API: getAllProducts and createProduct handlers use the gRPC client to interact with the server, convert responses to native structs, and return JSON.

- Background Streaming: StreamNewProducts runs in a background goroutine, logging each product received through streaming.

How to run: 
- Generate protobuff

  ```
  git clone https://github.com/googleapis/googleapis.git

  protoc -I . \
  -I /path/to/googleapis \
  --go_out gen --go_opt paths=source_relative \
  --go-grpc_out gen --go-grpc_opt paths=source_relative,require_unimplemented_servers=false \
  --grpc-gateway_out gen --grpc-gateway_opt paths=source_relative \
  product.proto

  ```

```
go run server/server.go server/main.go

go run gateway/main.go

go run client/main.go
```
