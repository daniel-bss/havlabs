package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/auth/token"
	"github.com/daniel-bss/havlabs/internal/auth/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RenewToken(ctx context.Context, req *pb.RenewTokenRequest) (*pb.RenewTokenResponse, error) {
	fmt.Println("OAKOAK")
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken, token.TokenTypeRefreshToken)
	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "session doesn't exist")
		}

		return nil, utils.InternalServerError()
	}

	if session.IsBlocked {
		return nil, status.Errorf(codes.NotFound, "this session is blocked")
	}

	if session.Username != refreshPayload.Username {
		return nil, status.Errorf(codes.NotFound, "usernames from token and session are mismatched")
	}

	if session.RefreshToken != req.RefreshToken {
		return nil, status.Errorf(codes.NotFound, "token and session payload are mismatched")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, status.Errorf(codes.NotFound, "refresh token is expired")
	}

	// CREATE TOKEN
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		refreshPayload.Role,
		server.config.AccessTokenDuration,
		token.TokenTypeAccessToken,
	)
	if err != nil {
		return nil, utils.FailedCreateAccessToken()
	}

	return &pb.RenewTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(time.Unix(int64(accessPayload.ExpiredAt), 0)),
	}, nil
}
