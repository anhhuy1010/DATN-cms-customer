package models

import (
	"context"
	"log"
	"time"

	"github.com/anhhuy1010/DATN-cms-customer/database"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	//"go.mongodb.org/mongo-driver/bson"

	"github.com/anhhuy1010/DATN-cms-customer/constant"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Customer struct {
	Uuid     string     `json:"uuid,omitempty" bson:"uuid"`
	UserName string     `json:"username" bson:"username"`
	Password string     `json:"password" bson:"password"`
	Email    string     `json:"email,omitempty" bson:"email"`
	IsActive int        `json:"is_active" bson:"is_active"`
	IsDelete int        `json:"is_delete" bson:"is_delete"`
	StartDay *time.Time `json:"startday" bson:"startday"`
	EndDay   *time.Time `json:"endday" bson:"endday"`
}

func (u *Customer) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("customer")
}

func (u *Customer) Find(conditions map[string]interface{}, opts ...*options.FindOptions) ([]*Customer, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	cursor, err := coll.Find(context.TODO(), conditions, opts...)
	if err != nil {
		return nil, err
	}

	var users []*Customer
	for cursor.Next(context.TODO()) {
		var elem Customer
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

func (u *Customer) Pagination(ctx context.Context, conditions map[string]interface{}, modelOptions ...ModelOption) ([]*Customer, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	modelOpt := ModelOption{}
	findOptions := modelOpt.GetOption(modelOptions)
	cursor, err := coll.Find(context.TODO(), conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var users []*Customer
	for cursor.Next(context.TODO()) {
		var elem Customer
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

func (u *Customer) Distinct(conditions map[string]interface{}, fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	values, err := coll.Distinct(context.TODO(), fieldName, conditions, opts...)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (u *Customer) FindOne(conditions map[string]interface{}) (*Customer, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	err := coll.FindOne(context.TODO(), conditions).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *Customer) Insert() (interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertOne(context.TODO(), u)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (u *Customer) InsertMany(Users []interface{}) ([]interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertMany(context.TODO(), Users)
	if err != nil {
		return nil, err
	}

	return resp.InsertedIDs, nil
}

func (u *Customer) Update() (int64, error) {
	coll := u.Model()

	condition := make(map[string]interface{})
	condition["uuid"] = u.Uuid

	updateStr := make(map[string]interface{})
	updateStr["$set"] = u

	resp, err := coll.UpdateOne(context.TODO(), condition, updateStr)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *Customer) UpdateByCondition(condition map[string]interface{}, data map[string]interface{}) (int64, error) {
	coll := u.Model()

	resp, err := coll.UpdateOne(context.TODO(), condition, data)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *Customer) UpdateMany(conditions map[string]interface{}, updateData map[string]interface{}) (int64, error) {
	coll := u.Model()
	resp, err := coll.UpdateMany(context.TODO(), conditions, updateData)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *Customer) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	coll := u.Model()

	condition["is_delete"] = constant.UNDELETE

	total, err := coll.CountDocuments(ctx, condition)
	if err != nil {
		return 0, err
	}

	return total, nil
}
func (u *Customer) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
