package user

type (
	GetListRequest struct {
		Keyword   string  `form:"keyword"`
		IdeasName *string `form:"ideasname"`
		Page      int     `form:"page"`
		Limit     int     `form:"limit"`
		Sort      string  `form:"sort"`
		IsActive  *int    `form:"is_active" `
		Role      *string `form:"role"`
	}
	ListResponse struct {
		Uuid          string `json:"uuid" `
		IsActive      int    `json:"is_active"`
		IdeasName     string `json:"ideasname"`
		Industry      string `json:"industry"`
		ContentDetail string `json:"content_detail"`
		Procedure     string `json:"procedure"`
		Price         int    `json:"price"`
		CustomerName  string `json:"customer_name"`
	}
)
