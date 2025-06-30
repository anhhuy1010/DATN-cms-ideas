package ideas

type (
	AddFavoriteRequest struct {
		PostUuid string `json:"post_uuid" binding:"required"`
	}
)
