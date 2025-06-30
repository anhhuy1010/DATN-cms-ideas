package ideas

import "time"

type (
	GetListRequestFavorite struct {
		Page     int    `form:"page"`
		Limit    int    `form:"limit"`
		Sort     string `form:"sort"`
		PostType string `form:"post_type"`
	}
	ListResponseFavorite struct {
		Uuid          string    `json:"uuid"`
		PostUuid      string    `json:"post_uuid"`
		CreatedAt     time.Time `json:"created_at"`
		CustomerUuid  string    `bson:"customer_uuid" json:"customer_uuid"`
		IdeasName     string    `json:"ideasname" bson:"ideasname"`
		Industry      string    `json:"industry" bson:"industry"`
		Procedure     string    `json:"is_procedure,omitempty" bson:"is_procedure"`
		ContentDetail string    `json:"content_detail,omitempty" bson:"content_detail"`
		View          int       `json:"view " bson:"view"`
		Price         int       `json:"price,omitempty" bson:"price"`
		PostDay       time.Time `json:"post_day" bson:"post_day"`
		CustomerName  string    `json:"customer_name" bson:"customer_name"`
		Image         string    `json:"image" bson:"image"`
	}
)
