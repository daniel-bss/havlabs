package dtos

import (
	db "github.com/daniel-bss/havlabs/internal/auth/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/auth/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
