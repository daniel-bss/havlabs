package utils

import (
	"context"
	"time"

	"google.golang.org/grpc/metadata"
)

func ExtractMetadataFromContextWithKey(key string, ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(key); len(values) > 0 {
			return values[0], nil
		}
	}

	return "", NewBadRequestError("invalid metadata")
}

func TimeFromString(s string, format string) time.Time {
	date, _ := time.Parse(format, s)
	return date
}
