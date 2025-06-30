package ideas

type (
	AddSoldIdeasRequest struct {
		PostUuid string `json:"post_uuid" binding:"required"`
		Status   int    `json:"status"`
	}
)
