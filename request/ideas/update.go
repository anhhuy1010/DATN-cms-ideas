package ideas

type (
	UpdateUri struct {
		Uuid string `uri:"uuid"`
	}
	UpdateRequest struct {
		IdeasName       string   `json:"ideasname"`
		Industry        string   `json:"industry"`
		ContentDetail   string   `json:"content_detail"`
		Value_Benefits  string   `json:"value_benefits"`
		Price           int64    `json:"price"`
		Image           []string `json:"image"`
		Image_Intellect *string  `json:"image_intellect"`
		Is_Intellect    int32    `json:"is_intellect"`
	}
)
