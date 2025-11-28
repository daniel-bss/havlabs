package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs/auth/pb"
	"github.com/daniel-bss/havlabs/auth/token"
	"github.com/daniel-bss/havlabs/auth/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

func (server *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	// user, err := server.store.GetUser(ctx, req.GetUsername())
	// if err != nil {
	// 	if errors.Is(err, db.ErrRecordNotFound) {
	// 		return nil, status.Errorf(codes.NotFound, "user not found")
	// 	}
	// 	return nil, status.Errorf(codes.Internal, "failed to find user")
	// }

	// err = util.CheckPassword(req.Password, user.HashedPassword)
	// if err != nil {
	// 	return nil, status.Errorf(codes.NotFound, "incorrect password")
	// }
	fmt.Println("MASOOK")

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		"myusername",
		"myrole",
		server.config.AccessTokenDuration,
		token.TokenTypeAccessToken,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		"myusername",
		"myrole",
		server.config.RefreshTokenDuration,
		token.TokenTypeRefreshToken,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	// mtdt := server.extractMetadata(ctx)

	// session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
	// 	ID:           refreshPayload.ID,
	// 	Username:     "heheusername",
	// 	RefreshToken: refreshToken,
	// 	UserAgent:    mtdt.UserAgent,
	// 	ClientIp:     mtdt.ClientIP,
	// 	IsBlocked:    false,
	// 	ExpiresAt:    refreshPayload.ExpiredAt,
	// })
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to create session")
	// }

	loginData := &pb.LoginData{
		// User:                  convertUser(user),
		// User: &pb.User{
		// 	Username:          "asd",
		// 	FullName:          "asd",
		// 	Email:             "asdas",
		// 	PasswordChangedAt: timestamppb.New(time.Now()),
		// 	CreatedAt:         timestamppb.New(time.Now()),
		// },
		User:                  &pb.User{},
		SessionId:             "sesisonid",
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	// Wrap the user data in google.protobuf.Any

	rsp := &pb.LoginResponse{
		Data: loginData,
	}

	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	if err := utils.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, utils.FieldViolation("password", err))
	}

	return violations
}
