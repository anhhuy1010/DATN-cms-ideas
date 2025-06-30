package ideas

type (
	GetPricelUri struct {
		Uuid string `uri:"uuid"`
	}
	GetPriceResponse struct {
		Price int `json:"price,omitempty"`
	}
)
