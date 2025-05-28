package service

import (
	"context"
	"errors"

	pb "github.com/anhhuy1010/DATN-cms-customer/grpc/proto/customer"
)

// CustomerServer implements the Customer gRPC service
type CustomerServer struct {
	pb.UnimplementedCustomerServer
	// Giả sử bạn có một biến db hoặc repository để truy xuất database ở đây
	customers map[string]*pb.CustomerResponse // tạm lưu trong bộ nhớ (demo)
}

func NewCustomerServer() *CustomerServer {
	return &CustomerServer{
		customers: make(map[string]*pb.CustomerResponse),
	}
}

// GetCustomer trả về thông tin khách hàng theo ID
func (s *CustomerServer) GetCustomer(ctx context.Context, req *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	customer, ok := s.customers[req.Uuid]
	if !ok {
		return nil, errors.New("customer not found")
	}
	return customer, nil
}
