package ideas

type (
	UpdateStatusUriAdmin struct {
		Uuid string `uri:"uuid"`
	}
	UpdateStatusRequestAdmin struct {
		IsDelete *int `json:"is_delete" binding:"omitempty,oneof=0 1"`
	}
)
