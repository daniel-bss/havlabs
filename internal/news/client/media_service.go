package client

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type MediaClient interface {
	GetMediaById(context.Context, uuid.UUID) (*uuid.UUID, error)
}

type mediaClientImpl struct {
	grpcClient pb.HavlabsMediaClient
}

func NewMediaClient(grpcConn *grpc.ClientConn) MediaClient {
	return &mediaClientImpl{
		grpcClient: pb.NewHavlabsMediaClient(grpcConn),
	}
}

func (mc *mediaClientImpl) GetMediaById(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	res, err := mc.grpcClient.GetMediaById(ctx, &pb.GetOneMediaByIdRequest{Id: id.String()})
	if err != nil {
		return nil, err
	}

	mediaId, err := uuid.Parse(res.Id)
	if err != nil {
		return nil, err
	}

	return &mediaId, nil
}
