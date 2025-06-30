package models

import (
	"context"
	"log"
	"time"

	"github.com/anhhuy1010/DATN-cms-ideas/constant"
	"github.com/anhhuy1010/DATN-cms-ideas/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Favorite struct {
	Uuid            string    `bson:"uuid" json:"uuid"`
	CustomerUuid    string    `bson:"customer_uuid" json:"customer_uuid"`
	PostUuid        string    `bson:"post_uuid" json:"post_uuid"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	IsDelete        int       `json:"is_delete" bson:"is_delete"`
	IdeasName       string    `json:"ideasname" bson:"ideasname"`
	Industry        string    `json:"industry" bson:"industry"`
	Procedure       string    `json:"is_procedure,omitempty" bson:"is_procedure"`
	ContentDetail   string    `json:"content_detail,omitempty" bson:"content_detail"`
	Value_Benefits  string    `json:"value_benefits,omitempty" bson:"value_benefits"`
	View            int       `json:"view " bson:"view"`
	Is_Intellect    int       `json:"is_intellect,omitempty" bson:"is_intellect"`
	Price           int       `json:"price,omitempty" bson:"price"`
	IsActive        int       `json:"is_active" bson:"is_active"`
	PostDay         time.Time `json:"post_day" bson:"post_day"`
	CustomerName    string    `json:"customer_name" bson:"customer_name"`
	CustomerEmail   string    `json:"customer_email" bson:"customer_email"`
	Image           []string  `json:"image" bson:"image"`
	Image_Intellect string    `json:"image_intellect" bson:"image_intellect"`
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
func (f *Favorite) Find(conditions map[string]interface{}, opts ...*options.FindOptions) ([]*Favorite, error) {
	coll := f.Model()
	conditions["is_delete"] = constant.UNDELETE
	cursor, err := coll.Find(context.TODO(), conditions, opts...)
	if err != nil {
		return nil, err
	}

	var users []*Favorite
	for cursor.Next(context.TODO()) {
		var elem Favorite
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
func (f *Favorite) Pagination(ctx context.Context, conditions map[string]interface{}, modelOptions ...ModelOption) ([]*Favorite, error) {
	coll := f.Model()
	conditions["is_delete"] = constant.UNDELETE
	modelOpt := ModelOption{}
	findOptions := modelOpt.GetOption(modelOptions)
	cursor, err := coll.Find(context.TODO(), conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var users []*Favorite
	for cursor.Next(context.TODO()) {
		var elem Favorite
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
func (f *Favorite) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	coll := f.Model()
	condition["is_delete"] = constant.UNDELETE
	total, err := coll.CountDocuments(ctx, condition)
	if err != nil {
		return 0, err
	}

	return total, nil
}
func (f *Favorite) FindOne(conditions map[string]interface{}) (*Favorite, error) {
	coll := f.Model()
	conditions["is_delete"] = constant.UNDELETE

	err := coll.FindOne(context.TODO(), conditions).Decode(&f)
	if err != nil {
		return nil, err
	}

	return f, nil
}
func (f *Favorite) Update() (int64, error) {
	coll := f.Model()
	condition := make(map[string]interface{})
	condition["uuid"] = f.Uuid

	updateStr := make(map[string]interface{})
	updateStr["$set"] = f

	resp, err := coll.UpdateOne(context.TODO(), condition, updateStr)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}
