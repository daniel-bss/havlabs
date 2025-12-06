package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/daniel-bss/havlabs-proto/pb"
	db "github.com/daniel-bss/havlabs/internal/auth/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/auth/dtos"
	"github.com/daniel-bss/havlabs/internal/auth/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AdminRole = "admin"
	UserRole  = "user"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{AdminRole, UserRole})
	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req, server.config)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	if authPayload.Role != AdminRole && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	fmt.Println(">>>", arg.FullName)

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: dtos.ConvertUser(user),
	}

	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest, config utils.Config) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateUsername(req.GetUsername(), config.MinUsernameLength, config.MaxUsernameLength); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	if req.Password != nil {
		if err := utils.ValidatePassword(req.GetPassword(), config.MinPwdLength, config.MaxPwdLength); err != nil {
			violations = append(violations, utils.FieldViolation("password", err))
		}
	}

	if req.FullName != nil {
		if err := utils.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, utils.FieldViolation("full_name", err))
		}
	}

	// if req.Email != nil {
	// 	if err := utils.ValidateEmail(req.GetEmail()); err != nil {
	// 		violations = append(violations, utils.FieldViolation("email", err))
	// 	}
	// }

	return violations
}
