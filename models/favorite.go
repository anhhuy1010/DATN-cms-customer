package models

import (
	"context"
	"time"

	"github.com/anhhuy1010/DATN-cms-customer/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type Favorite struct {
	Uuid         string    `bson:"uuid" json:"uuid"`
	CustomerUuid string    `bson:"customer_uuid" json:"customer_uuid"`
	PostUuid     string    `bson:"post_uuid" json:"post_uuid"`
	PostType     string    `bson:"post_type" json:"post_type"` // "idea" | "problem"
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

func (f *Favorite) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("favorites")
}

func (f *Favorite) Insert() (interface{}, error) {
	coll := f.Model()

	resp, err := coll.InsertOne(context.TODO(), f)
	if err != nil {
		return 0, err
	}

	return resp, nil
}
