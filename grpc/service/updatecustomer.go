package service

import (
	"context"
	"time"

	pb "github.com/anhhuy1010/DATN-cms-customer/grpc/proto/updatecustomer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateCustomerServer struct {
	pb.UnimplementedUpdateCustomerServer
	Collection *mongo.Collection // MongoDB collection
}

// NewUpdateCustomerServer creates a new UpdateCustomerServer with the given MongoDB collection.
func NewUpdateCustomerServer(collection *mongo.Collection) *UpdateCustomerServer {
	return &UpdateCustomerServer{
		Collection: collection,
	}
}

func (u *UpdateCustomerServer) UpdateDuration(ctx context.Context, req *pb.UpdateDurationRequest) (*pb.UpdateDurationResponse, error) {
	filter := bson.M{"uuid": req.Uuid}
	now := time.Now()
	end := now.AddDate(0, 0, 30)

	update := bson.M{
		"$set": bson.M{
			"startday": now,
			"endday":   end,
		},
	}

	_, err := u.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot update duration: %v", err)
	}

	return &pb.UpdateDurationResponse{
		Message: "Customer duration updated successfully",
	}, nil
}
