package models

import (
	"context"
	"log"
	"time"

	"github.com/anhhuy1010/DATN-cms-ideas/database"
	"go.mongodb.org/mongo-driver/mongo"

	//"go.mongodb.org/mongo-driver/bson"

	"github.com/anhhuy1010/DATN-cms-ideas/constant"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Ideas struct {
	Uuid           string    `json:"uuid,omitempty" bson:"uuid"`
	CustomerUuid   string    `json:"customeruuid,omitempty" bson:"customeruuid"`
	IdeasName      string    `json:"ideasname" bson:"ideasname"`
	Industry       string    `json:"industry" bson:"industry"`
	OrtherIndustry string    `json:"orderindustry,omitempty" bson:"orderindustry"`
	IsProcedure    int       `json:"is_procedure,omitempty" bson:"is_procedure"`
	ContentDetail  string    `json:"content_detail,omitempty" bson:"content_detail"`
	Value_Benefits string    `json:"value_benefits,omitempty" bson:"value_benefits"`
	Is_Intellect   int       `json:"is_intellect,omitempty" bson:"is_intellect"`
	Price          int       `json:"price,omitempty" bson:"price"`
	IsActive       int       `json:"is_active" bson:"is_active"`
	IsDelete       int       `json:"is_delete" bson:"is_delete"`
	PostDay        time.Time `json:"post_day" bson:"post_day"`
	CustomerName   string    `json:"customer_name" bson:"customer_name"`
	CustomerEmail  string    `json:"customer_email" bson:"customer_email"`
}

func (u *Ideas) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("ideas")
}

func (u *Ideas) Find(conditions map[string]interface{}, opts ...*options.FindOptions) ([]*Ideas, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	cursor, err := coll.Find(context.TODO(), conditions, opts...)
	if err != nil {
		return nil, err
	}

	var users []*Ideas
	for cursor.Next(context.TODO()) {
		var elem Ideas
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

func (u *Ideas) Pagination(ctx context.Context, conditions map[string]interface{}, modelOptions ...ModelOption) ([]*Ideas, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	modelOpt := ModelOption{}
	findOptions := modelOpt.GetOption(modelOptions)
	cursor, err := coll.Find(context.TODO(), conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var users []*Ideas
	for cursor.Next(context.TODO()) {
		var elem Ideas
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

func (u *Ideas) Distinct(conditions map[string]interface{}, fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	values, err := coll.Distinct(context.TODO(), fieldName, conditions, opts...)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (u *Ideas) FindOne(conditions map[string]interface{}) (*Ideas, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	err := coll.FindOne(context.TODO(), conditions).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *Ideas) Insert(ctx context.Context) error {
	coll := u.Model()
	_, err := coll.InsertOne(ctx, u)
	return err
}

func (u *Ideas) InsertMany(Users []interface{}) ([]interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertMany(context.TODO(), Users)
	if err != nil {
		return nil, err
	}

	return resp.InsertedIDs, nil
}

func (u *Ideas) Update() (int64, error) {
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

func (u *Ideas) UpdateByCondition(condition map[string]interface{}, data map[string]interface{}) (int64, error) {
	coll := u.Model()

	resp, err := coll.UpdateOne(context.TODO(), condition, data)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *Ideas) UpdateMany(conditions map[string]interface{}, updateData map[string]interface{}) (int64, error) {
	coll := u.Model()
	resp, err := coll.UpdateMany(context.TODO(), conditions, updateData)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *Ideas) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	coll := u.Model()

	condition["is_delete"] = constant.UNDELETE

	total, err := coll.CountDocuments(ctx, condition)
	if err != nil {
		return 0, err
	}

	return total, nil
}
