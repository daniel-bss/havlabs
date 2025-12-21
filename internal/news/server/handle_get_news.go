package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/news/dtos"
	"github.com/daniel-bss/havlabs/internal/news/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) GetPaginatedNews(ctx context.Context, req *pb.ListNewsRequest) (*pb.ListNewsResponse, error) {
	news, tot, err := server.uc.GetNews(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("error from media/GetNews")

		if e, ok := err.(utils.BadRequestError); ok {
			return nil, status.Error(codes.InvalidArgument, e.Error())
		}
		return nil, err
	}

	result := utils.Map(news, func(b dtos.NewsDto) *pb.OneNewsResponse {
		return &pb.OneNewsResponse{
			Title:       b.Title,
			Content:     b.Content,
			ImageUrl:    b.ImageURL,
			PublishedAt: timestamppb.New(b.PublishedAt.Time),
		}
	})

	return &pb.ListNewsResponse{
		News:       result,
		TotalPages: tot,
	}, nil
}

func (server *Server) GetOneNews(ctx context.Context, req *pb.GetOneNewsByIdRequest) (*pb.OneNewsResponse, error) {
	fmt.Println("get one")
	return &pb.OneNewsResponse{Title: "TEST"}, nil
}
