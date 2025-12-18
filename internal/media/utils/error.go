package utils

import (
	"github.com/jackc/pgx/v5"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrRecordNotFound = pgx.ErrNoRows

func FieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func InvalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")

	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}

func UnauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}

func InternalServerError() error {
	return status.Errorf(codes.Internal, "UNEXPECTED_ERROR")
}

func FailedCreateAccessToken() error {
	return status.Errorf(codes.Internal, "failed to create access token")
}

type BadRequestError struct {
	msg string
}

func NewBadRequestError(s string) error {
	return BadRequestError{msg: s}
}

func (e BadRequestError) Error() string {
	return e.msg
}
