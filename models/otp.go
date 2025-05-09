package models

import (
	"context"
	"time"

	"github.com/anhhuy1010/cms-user/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OTP struct {
	Uuid      string    `bson:"uuid"`
	UserUuid  string    `bson:"user_uuid"`
	OtpCode   string    `bson:"otp_code"`
	ExpiresAt time.Time `bson:"expires_at"`
	CreatedAt time.Time `bson:"created_at"`
}

func (o *OTP) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("otp")
}

// Lưu OTP mới
func (o *OTP) Insert(ctx context.Context) error {
	coll := o.Model()
	_, err := coll.InsertOne(ctx, o)
	return err
}

// Tìm OTP theo Email
func (o *OTP) FindOTPByEmail(ctx context.Context, email string) (*OTP, error) {
	col := o.Model()
	filter := bson.M{"email": email}
	var otp OTP
	err := col.FindOne(ctx, filter).Decode(&otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// Xoá OTP theo Email
func (o *OTP) DeleteOTP(ctx context.Context, email string) error {
	col := o.Model()
	filter := bson.M{"email": email}
	_, err := col.DeleteOne(ctx, filter)
	return err
}

// Kiểm tra OTP có hết hạn không (ví dụ timeout 5 phút)
func (o *OTP) IsExpired() bool {
	return time.Since(o.CreatedAt) > 5*time.Minute
}
