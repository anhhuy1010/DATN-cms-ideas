package service

import (
	"context"
	"time"

	idea "github.com/anhhuy1010/DATN-cms-ideas/grpc/proto/idea"
	"github.com/anhhuy1010/DATN-cms-ideas/models"
	"github.com/google/uuid"
)

type IdeaServiceServer struct {
	idea.UnimplementedIdeaServiceServer
}

func NewIdeaServiceServer() idea.IdeaServiceServer {
	return &IdeaServiceServer{}
}

func (s *IdeaServiceServer) CreateIdea(ctx context.Context, req *idea.CreateIdeaRequest) (*idea.CreateIdeaResponse, error) {
	idModel := models.Ideas{
		Uuid:           uuid.New().String(),
		IdeasName:      req.IdeasName,
		Industry:       req.Industry,
		ContentDetail:  req.ContentDetail,
		Price:          int(req.Price),
		PostDay:        time.Now(),
		CustomerUuid:   req.CustomerUuid,
		CustomerName:   req.CustomerName,
		CustomerEmail:  req.CustomerEmail,
		IsActive:       1,
		IsDelete:       0,
		Is_Intellect:   0,
		Value_Benefits: "",
		Procedure:      "",
		Image:          []string{},
	}

	err := idModel.Insert(ctx)
	if err != nil {
		return nil, err
	}

	return &idea.CreateIdeaResponse{
		Uuid:    idModel.Uuid,
		Message: "Ý tưởng đã được tạo thành công",
	}, nil
}
