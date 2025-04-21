package service

import (
	"context"

	idea "github.com/anhhuy1010/DATN-cms-ideas/grpc/proto/idea"
	"github.com/anhhuy1010/DATN-cms-ideas/models"
)

type IdeaServiceServer struct {
}

func (s *IdeaServiceServer) ListIdeas(ctx context.Context, req *idea.ListIdeasRequest) (*idea.ListIdeasResponse, error) {
	ideaModel := models.Ideas{}

	// Tính toán phân trang
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}
	skip := int64((page - 1) * limit)

	// Tạo ModelOption với sắp xếp theo post_day giảm dần
	opt := models.ModelOption{
		SortBy:  "post_day",
		SortDir: -1,
		Limit:   int64(limit),
		Skip:    skip,
	}

	ideas, err := ideaModel.Pagination(ctx, map[string]interface{}{}, opt)
	if err != nil {
		return nil, err
	}

	var result []*idea.Idea
	for _, i := range ideas {
		result = append(result, &idea.Idea{
			Uuid:          i.Uuid,
			Ideasname:     i.IdeasName,
			Industry:      i.Industry,
			ContentDetail: i.ContentDetail,
			Price:         int32(i.Price),
			PostDay:       i.PostDay.Format("07-02-1999"),
		})
	}

	return &idea.ListIdeasResponse{Ideas: result}, nil
}
