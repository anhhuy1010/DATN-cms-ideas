package ideas

type (
	GetListRequest struct {
		Keyword   string  `form:"keyword"`
		Ideasname *string `form:"ideasname"`
		Page      int     `form:"page"`
		Limit     int     `form:"limit"`
		Sort      string  `form:"sort"`
		IsActive  *int    `form:"is_active" `
		Role      *string `form:"role"`
	}
	ListResponse struct {
		Uuid      string `json:"uuid" `
		IdeasName string `json:"username"`
		IsActive  int    `json:"is_active"`
		Industry  string `json:"industry"`
	}
)
