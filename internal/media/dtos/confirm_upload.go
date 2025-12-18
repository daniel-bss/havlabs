package dtos

import "github.com/daniel-bss/havlabs-proto/pb"

type ConfirmUpload struct {
	MediaId string
	Status  pb.StatusEnum
}
