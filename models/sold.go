package models

import (
	"context"
	"log"
	"time"

	"github.com/anhhuy1010/DATN-cms-ideas/database"
	"go.mongodb.org/mongo-driver/mongo"

	//"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Sold struct {
	Uuid         string    `json:"uuid,omitempty" bson:"uuid"`
	PostUuid     string    `bson:"post_uuid" json:"post_uuid"`
	CustomerUuid string    `json:"customeruuid,omitempty" bson:"customeruuid"`
	CreatedAt    time.Time `json:"created_at,omitempty" bson:"created_at"`
}

func (s *Sold) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("sold")
}

func (s *Sold) Find(conditions map[string]interface{}, opts ...*options.FindOptions) ([]*Sold, error) {
	coll := s.Model()

	cursor, err := coll.Find(context.TODO(), conditions, opts...)
	if err != nil {
		return nil, err
	}

	var users []*Sold
	for cursor.Next(context.TODO()) {
		var elem Sold
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

// //////////////////////////////////////////////////////////////////////////////////////////////
func (s *Sold) Pagination(ctx context.Context, conditions map[string]interface{}, modelOptions ...ModelOption) ([]*Sold, error) {
	coll := s.Model()

	modelOpt := ModelOption{}
	findOptions := modelOpt.GetOption(modelOptions)
	cursor, err := coll.Find(ctx, conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var users []*Sold
	for cursor.Next(ctx) {
		var elem Sold
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

// ///////////////////////////////////////////////////////////////////////////////////////////////
func (s *Sold) Distinct(conditions map[string]interface{}, fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	coll := s.Model()

	values, err := coll.Distinct(context.TODO(), fieldName, conditions, opts...)
	if err != nil {
		return nil, err
	}

	return values, nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
func (s *Sold) FindOne(conditions map[string]interface{}) (*Sold, error) {
	coll := s.Model()
	err := coll.FindOne(context.TODO(), conditions).Decode(&s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// ///////////////////////////////////////////////////////////////////////////////////////////
func (s *Sold) Insert(ctx context.Context) error {
	coll := s.Model()
	_, err := coll.InsertOne(ctx, s)
	return err
}

// //////////////////////////////////////////////////////////
func (s *Sold) InsertMany(Users []interface{}) ([]interface{}, error) {
	coll := s.Model()

	resp, err := coll.InsertMany(context.TODO(), Users)
	if err != nil {
		return nil, err
	}

	return resp.InsertedIDs, nil
}

// ///////////////////////////////////////////////////////////////////////////////////
func (s *Sold) Update() (int64, error) {
	coll := s.Model()

	condition := make(map[string]interface{})
	condition["uuid"] = s.Uuid

	updateStr := make(map[string]interface{})
	updateStr["$set"] = s

	resp, err := coll.UpdateOne(context.TODO(), condition, updateStr)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

// ///////////////////////////////////////////////////////////////////////////////////
func (s *Sold) UpdateByCondition(condition map[string]interface{}, data map[string]interface{}) (int64, error) {
	coll := s.Model()

	resp, err := coll.UpdateOne(context.TODO(), condition, data)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

// ////////////////////////////////////////////////////////////////////////////////////
func (s *Sold) UpdateMany(conditions map[string]interface{}, updateData map[string]interface{}) (int64, error) {
	coll := s.Model()
	resp, err := coll.UpdateMany(context.TODO(), conditions, updateData)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

// //////////////////////////////////////////////////////////////////////////////////////
func (s *Sold) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	coll := s.Model()

	total, err := coll.CountDocuments(ctx, condition)
	if err != nil {
		return 0, err
	}

	return total, nil
}
