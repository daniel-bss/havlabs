package server

import (
	"context"

	"github.com/daniel-bss/havlabs/auth/pb"
	"github.com/daniel-bss/havlabs/auth/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/anypb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.BaseResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	// hashedPassword, err := utils.HashPassword(req.GetPassword())
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	// }

	// arg := db.CreateUserTxParams{
	// 	CreateUserParams: db.CreateUserParams{
	// 		Username:       req.GetUsername(),
	// 		HashedPassword: hashedPassword,
	// 		FullName:       req.GetFullName(),
	// 		Email:          req.GetEmail(),
	// 	},
	// 	AfterCreate: func(user db.User) error {
	// 		taskPayload := &worker.PayloadSendVerifyEmail{
	// 			Username: user.Username,
	// 		}
	// 		opts := []asynq.Option{
	// 			asynq.MaxRetry(10),
	// 			asynq.ProcessIn(10 * time.Second),
	// 			asynq.Queue(worker.QueueCritical),
	// 		}

	// 		return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
	// 	},
	// }

	// txResult, err := server.store.CreateUserTx(ctx, arg)
	// if err != nil {
	// 	if db.ErrorCode(err) == db.UniqueViolation {
	// 		return nil, status.Error(codes.AlreadyExists, err.Error())
	// 	}
	// 	return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	// }

	createUserrsp := &pb.CreateUserResponse{
		// User: convertUser(txResult.User),
		User: &pb.User{},
	}

	anyData, err := anypb.New(createUserrsp)
	if err != nil {
		return nil, err
	}

	rsp := &pb.BaseResponse{
		Data: anyData,
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	if err := utils.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, utils.FieldViolation("password", err))
	}

	if err := utils.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, utils.FieldViolation("full_name", err))
	}

	if err := utils.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, utils.FieldViolation("email", err))
	}

	return violations
}
