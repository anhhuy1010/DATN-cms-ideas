package ideas

import "time"

type (
	GetListWattingRequest struct {
		PriceTier *string `form:"price_tier"`
		Keyword   string  `form:"keyword"`
		Ideasname *string `form:"ideasname"`
		Industry  *string `form:"industry"`
		View      *int    `form:"view"`
		Page      int     `form:"page"`
		Limit     int     `form:"limit"`
		Sort      string  `form:"sort"`
		IsActive  *int    `form:"is_active" `
		IsDelete  *int    `form:"is_delete" `
		Price     *int    `form:"price"`
		Uuid      *string `form:"uuid"`
	}
	ListWattingResponse struct {
		Uuid         string    `json:"uuid" `
		CustomerName string    `json:"customer_name"`
		IdeasName    string    `json:"ideasname"`
		Industry     string    `json:"industry"`
		IsActive     int       `json:"is_active"`
		IsDelete     int       `json:"is_delete"`
		PostDay      time.Time `json:"post_day"`
	}
)
