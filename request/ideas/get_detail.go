package ideas

import "time"

type (
	GetDetailUri struct {
		Uuid string `uri:"uuid"`
	}
	GetDetailResponse struct {
		IdeasName       string    `json:"ideasname"`
		Industry        string    `json:"industry"`
		OrtherIndustry  string    `json:"orderindustry,omitempty"`
		Procedure       string    `json:"is_procedure,omitempty"`
		ContentDetail   string    `json:"content_detail,omitempty"`
		Value_Benefits  string    `json:"value_benefits,omitempty"`
		Is_Intellect    int       `json:"is_intellect,omitempty"`
		Price           int       `json:"price,omitempty"`
		CustomerName    string    `json:"customer_name"`
		CustomerEmail   string    `json:"customer_email"`
		Image           []string  `json:"image"`
		View            int       `json:"view" bson:"view"`
		Image_Intellect string    `json:"image_intellect"`
		CustomerUuid    string    `json:"customeruuid"`
		PostDay         time.Time `json:"post_day"`
		IsActive        int       `json:"is_active"`
		IsDelete        int       `json:"is_delete"`
		Status          int       `json:"status"`
	}
)
