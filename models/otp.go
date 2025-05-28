package models

import (
	"context"
	"fmt"
	"time"

	"github.com/anhhuy1010/DATN-cms-customer/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OTP struct {
	Uuid      string    `json:"uuid" bson:"uuid"`
	UserUuid  string    `json:"user_uuid" bson:"user_uuid"`
	OtpCode   string    `json:"otp_code" bson:"otp_code"`
	Email     string    `json:"email" bson:"email"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
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

	// Kiểm tra thời gian hết hạn
	if time.Now().After(otp.ExpiresAt) {
		// Nếu đã hết hạn, xóa bản ghi OTP này
		_, _ = col.DeleteOne(ctx, filter)
		return nil, fmt.Errorf("OTP đã hết hạn")
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
	return time.Since(o.CreatedAt) > 30*time.Second
}
