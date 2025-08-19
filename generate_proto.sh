#!/bin/bash

set -e

# Đường dẫn tới thư mục proto
PROTO_DIR="./proto"

# Thư mục chứa output sinh ra
OUT_DIR="./proto/pb"

# Các thư viện phụ thuộc (đã clone)
GOOGLE_APIS="./proto/google/api"
GRPC_GATEWAY="./proto/protoc-gen-openapiv2/options"

# Tạo thư mục đầu ra nếu chưa có
mkdir -p ${OUT_DIR}

echo "🔧 Building proto files..."

# Duyệt tất cả các file .proto trong thư mục proto/
for file in ${PROTO_DIR}/*.proto; do
    echo "📦 Compiling $file"

    protoc \
      --proto_path=${PROTO_DIR} \
      --proto_path=${GOOGLE_APIS} \
      --proto_path=${GRPC_GATEWAY} \
      --go_out=${OUT_DIR} \
      --go_opt=paths=source_relative \
      --go-grpc_out=${OUT_DIR} \
      --go-grpc_opt=paths=source_relative \
      --grpc-gateway_out=${OUT_DIR} \
      --grpc-gateway_opt=paths=source_relative \
      --grpc-gateway_opt=generate_unbound_methods=true \
      --openapiv2_out=${OUT_DIR} \
      --openapiv2_opt=logtostderr=true \
      "$file"
done

echo "✅ Done generating Go code from proto!"
