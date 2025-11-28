package server

import (
	"context"

	"github.com/daniel-bss/havlabs/auth/pb"
	"github.com/daniel-bss/havlabs/auth/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	DepositorRole = "depositor"
	BankerRole    = "banker"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.BaseResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{BankerRole, DepositorRole})
	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	if authPayload.Role != BankerRole && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	// arg := db.UpdateUserParams{
	// 	Username: req.GetUsername(),
	// 	FullName: pgtype.Text{
	// 		String: req.GetFullName(),
	// 		Valid:  req.FullName != nil,
	// 	},
	// 	Email: pgtype.Text{
	// 		String: req.GetEmail(),
	// 		Valid:  req.Email != nil,
	// 	},
	// }

	// if req.Password != nil {
	// 	hashedPassword, err := util.HashPassword(req.GetPassword())
	// 	if err != nil {
	// 		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	// 	}

	// 	arg.HashedPassword = pgtype.Text{
	// 		String: hashedPassword,
	// 		Valid:  true,
	// 	}

	// 	arg.PasswordChangedAt = pgtype.Timestamptz{
	// 		Time:  time.Now(),
	// 		Valid: true,
	// 	}
	// }

	// user, err := server.store.UpdateUser(ctx, arg)
	// if err != nil {
	// 	if errors.Is(err, db.ErrRecordNotFound) {
	// 		return nil, status.Errorf(codes.NotFound, "user not found")
	// 	}
	// 	return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	// }

	updateUserRsp := &pb.UpdateUserResponse{
		User: &pb.User{},
	}

	anyData, err := anypb.New(updateUserRsp)
	if err != nil {
		return nil, err
	}

	rsp := &pb.BaseResponse{
		Data: anyData,
	}

	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	if req.Password != nil {
		if err := utils.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, utils.FieldViolation("password", err))
		}
	}

	if req.FullName != nil {
		if err := utils.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, utils.FieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := utils.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, utils.FieldViolation("email", err))
		}
	}

	return violations
}
