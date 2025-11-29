package server

import (
	"context"
	"fmt"

	db "github.com/daniel-bss/havlabs/auth/db/sqlc"
	"github.com/daniel-bss/havlabs/auth/dtos"
	"github.com/daniel-bss/havlabs/auth/pb"
	"github.com/daniel-bss/havlabs/auth/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req, server.config)
	if violations != nil {
		for _, v := range violations {
			fmt.Println(v)
		}
		return nil, utils.InvalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
		},
		AfterCreate: func(user db.User) error { return nil },
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: dtos.ConvertUser(txResult.User),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest, config utils.Config) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateUsername(req.GetUsername(), config.MinUsernameLength, config.MaxUsernameLength); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	if err := utils.ValidatePassword(req.GetPassword(), config.MinPwdLength, config.MaxPwdLength); err != nil {
		violations = append(violations, utils.FieldViolation("password", err))
	}

	if err := utils.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, utils.FieldViolation("full_name", err))
	}

	// TODO: EMAIL
	// if err := utils.ValidateEmail(req.GetEmail()); err != nil {
	// 	violations = append(violations, utils.FieldViolation("email", err))
	// }

	return violations
}
