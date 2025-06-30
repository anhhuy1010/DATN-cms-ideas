package ideas

type (
	UpdateStatusUris struct {
		Uuid string `uri:"uuid"`
	}
	RejectStatusRequest struct {
		IsDelete *int `json:"is_delete" binding:"required"`
	}
)
