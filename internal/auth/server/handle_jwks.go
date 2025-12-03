package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/daniel-bss/havlabs/internal/auth/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetJWKS(ctx context.Context, req *emptypb.Empty) (res *pb.JWKSResponse, err error) {
	fmt.Println("!@#!@#!@#!@#")
	pubKey := s.tokenMaker.PublicKey()

	jwk := &pb.JWK{
		Kid: "oneandonly-jwk-until-further-implementation", // OPTIONAL UNTIL MULTIPLE JWKs
		Kty: "RSA",
		Alg: "RS384",
		Use: "sig",
		N:   base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pubKey.E)).Bytes()),
	}

	return &pb.JWKSResponse{
		Keys: []*pb.JWK{jwk},
	}, nil
}
