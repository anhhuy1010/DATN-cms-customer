package models

import (
	"context"

	"github.com/anhhuy1010/DATN-cms-customer/database"
	"go.mongodb.org/mongo-driver/mongo"

	//"go.mongodb.org/mongo-driver/bson"

	"github.com/anhhuy1010/DATN-cms-customer/constant"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tokens struct {
	UserUuid string `json:"user_uuid,omitempty" bson:"user_uuid"`
	Uuid     string `json:"uuid,omitempty" bson:"uuid"`
	Token    string `json:"token"`
	IsDelete int    `json:"is_delete" bson:"is_delete"`
}

func (u *Tokens) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("tokens")
}

func (u *Tokens) Find(conditions map[string]interface{}, opts ...*options.FindOptions) ([]*Tokens, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	cursor, err := coll.Find(context.TODO(), conditions, opts...)
	if err != nil {
		return nil, err
	}

	var tokens []*Tokens
	for cursor.Next(context.TODO()) {
		var elem Tokens
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	_ = cursor.Close(context.TODO())

	return tokens, nil
}

func (u *Tokens) FindOne(conditions map[string]interface{}) (*Tokens, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	err := coll.FindOne(context.TODO(), conditions).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *Tokens) Insert() (interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertOne(context.TODO(), u)
	if err != nil {
		return 0, err
	}

	return resp, nil
}
