package models

import pb "product-service/proto/gen"

// Product struct for GORM and SQLite database
type Product struct {
	ID    string  `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

// ToProto converts a Product struct to a protobuf Product message
func (p *Product) ToProto() *pb.Product {
	return &pb.Product{
		Id:    p.ID,
		Name:  p.Name,
		Price: p.Price,
	}
}

// ProductFromProto converts a protobuf Product message to a Product struct
func ProductFromProto(proto *pb.Product) *Product {
	return &Product{
		ID:    proto.Id,
		Name:  proto.Name,
		Price: proto.Price,
	}
}
