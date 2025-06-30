package ideas

import "time"

type (
	GetListRequest struct {
		PriceTier   *string `form:"price_tier"`
		Keyword     string  `form:"keyword"`
		Ideasname   *string `form:"ideasname"`
		Industry    *string `form:"industry"`
		View        *int    `form:"view"`
		Page        int     `form:"page"`
		Limit       int     `form:"limit"`
		Sort        string  `form:"sort"`
		IsActive    *int    `form:"is_active" `
		IsDelete    *int    `form:"is_delete" `
		Price       *int    `form:"price"`
		TopViewOnly bool    `form:"top_view_only"`
	}
	ListResponse struct {
		Uuid          string    `json:"uuid" `
		CustomerName  string    `json:"customer_name"`
		IdeasName     string    `json:"ideasname"`
		Industry      string    `json:"industry"`
		Image         string    `json:"image" bson:"image"`
		ContentDetail string    `json:"content_detail" bson:"content_detail"`
		Price         int       `json:"price" bson:"price"`
		PostDay       time.Time `json:"post_day" bson:"post_day"`
		View          int       `json:"view" bson:"view"`
		IsActive      int       `json:"is_active" `
		IsDelete      int       `json:"is_delete" `
	}
)
