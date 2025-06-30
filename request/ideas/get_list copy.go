package ideas

import "time"

type (
	ListBuyResponse struct {
		Uuid          string    `json:"uuid" `
		CustomerName  string    `json:"customer_name"`
		IdeasName     string    `json:"ideasname"`
		Industry      string    `json:"industry"`
		Image         string    `json:"image" bson:"image"`
		ContentDetail string    `json:"content_detail" bson:"content_detail"`
		PostDay       time.Time `json:"post_day" bson:"post_day"`
	}
)
