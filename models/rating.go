package models

import (
	"context"
	"log"
	"time"

	"github.com/anhhuy1010/DATN-cms-customer/constant"
	"github.com/anhhuy1010/DATN-cms-customer/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Rating struct {
	Uuid         string    `json:"uuid,omitempty" bson:"uuid"`
	CustomerUuid string    `bson:"customer_uuid" json:"customer_uuid"`
	CustomerName string    `bson:"customer_name" json:"customer_name"`
	ExpertUuid   string    `bson:"expert_uuid" json:"expert_uuid"`
	Rating       int       `bson:"rating" json:"rating"` // 1 đến 5
	Comment      string    `bson:"comment" json:"comment"`
	IsDelete     int       `json:"is_delete" bson:"is_delete"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

func (r *Rating) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("rating")
}

func (r *Rating) Find(conditions map[string]interface{}, opts ...*options.FindOptions) ([]*Rating, error) {
	coll := r.Model()

	conditions["is_delete"] = constant.UNDELETE
	cursor, err := coll.Find(context.TODO(), conditions, opts...)
	if err != nil {
		return nil, err
	}

	var users []*Rating
	for cursor.Next(context.TODO()) {
		var elem Rating
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}

		users = append(users, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	_ = cursor.Close(context.TODO())

	return users, nil
}

func (r *Rating) Pagination(ctx context.Context, conditions map[string]interface{}, modelOptions ...ModelOption) ([]*Rating, error) {
	coll := r.Model()

	conditions["is_delete"] = constant.UNDELETE

	modelOpt := ModelOption{}
	findOptions := modelOpt.GetOption(modelOptions)
	cursor, err := coll.Find(context.TODO(), conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var users []*Rating
	for cursor.Next(context.TODO()) {
		var elem Rating
		err := cursor.Decode(&elem)
		if err != nil {
			log.Println("[Decode] PopularCuisine:", err)
			log.Println("-> #", elem.Uuid)
			continue
		}

		users = append(users, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	_ = cursor.Close(context.TODO())

	return users, nil
}

func (r *Rating) Distinct(conditions map[string]interface{}, fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	coll := r.Model()

	conditions["is_delete"] = constant.UNDELETE

	values, err := coll.Distinct(context.TODO(), fieldName, conditions, opts...)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (r *Rating) FindOne(conditions map[string]interface{}) (*Rating, error) {
	coll := r.Model()

	conditions["is_delete"] = constant.UNDELETE
	err := coll.FindOne(context.TODO(), conditions).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Rating) Insert() (interface{}, error) {
	coll := r.Model()

	resp, err := coll.InsertOne(context.TODO(), r)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (r *Rating) InsertMany(Users []interface{}) ([]interface{}, error) {
	coll := r.Model()

	resp, err := coll.InsertMany(context.TODO(), Users)
	if err != nil {
		return nil, err
	}

	return resp.InsertedIDs, nil
}

func (r *Rating) Update() (int64, error) {
	coll := r.Model()

	condition := make(map[string]interface{})
	condition["uuid"] = r.Uuid

	updateStr := make(map[string]interface{})
	updateStr["$set"] = r

	resp, err := coll.UpdateOne(context.TODO(), condition, updateStr)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (r *Rating) UpdateByCondition(condition map[string]interface{}, data map[string]interface{}) (int64, error) {
	coll := r.Model()

	resp, err := coll.UpdateOne(context.TODO(), condition, data)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (r *Rating) UpdateMany(conditions map[string]interface{}, updateData map[string]interface{}) (int64, error) {
	coll := r.Model()
	resp, err := coll.UpdateMany(context.TODO(), conditions, updateData)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (r *Rating) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	coll := r.Model()

	condition["is_delete"] = constant.UNDELETE

	total, err := coll.CountDocuments(ctx, condition)
	if err != nil {
		return 0, err
	}

	return total, nil
}
