package ideas

type (
	UpdateStatusUri struct {
		Uuid string `uri:"uuid"`
	}
	AcceptStatusRequest struct {
		IsActive *int `json:"is_active" binding:"required"`
	}
)
