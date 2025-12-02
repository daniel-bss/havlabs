package server

import (
	"context"
	"errors"
	"fmt"

	db "github.com/daniel-bss/havlabs/internal/auth/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/auth/dtos"
	"github.com/daniel-bss/havlabs/internal/auth/pb"
	"github.com/daniel-bss/havlabs/internal/auth/token"
	"github.com/daniel-bss/havlabs/internal/auth/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	violations := validateLoginUserRequest(req, server.config)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		string(user.Role),
		server.config.AccessTokenDuration,
		token.TokenTypeAccessToken,
	)
	if err != nil {
		fmt.Printf("failed to create access token: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		string(user.Role),
		server.config.RefreshTokenDuration,
		token.TokenTypeRefreshToken,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	mtdt := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	loginData := &pb.LoginData{
		User:                  dtos.ConvertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
	}

	rsp := &pb.LoginResponse{
		Data: loginData,
	}

	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginRequest, config utils.Config) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateUsername(req.GetUsername(), config.MinUsernameLength, config.MaxUsernameLength); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	if err := utils.ValidatePassword(req.GetPassword(), config.MinPwdLength, config.MaxPwdLength); err != nil {
		violations = append(violations, utils.FieldViolation("password", err))
	}

	return violations
}
