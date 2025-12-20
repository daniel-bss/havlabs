package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"google.golang.org/grpc/metadata"
)

var purposes = map[string]bool{
	"news": true,
}

var contentTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
}

const maxMIMESniffBytes = 512

func DetectMIME(r io.Reader) (string, error) {
	buf := make([]byte, maxMIMESniffBytes)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	return http.DetectContentType(buf[:n]), nil
}

func GetChecksum(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return hex.EncodeToString(h.Sum(nil))
}

func IsValidPurpose(s string) bool {
	_, ok := purposes[s]
	return ok
}

func IsValidContentType(s string) bool {
	if s == "image/webp" {
		return true
	}
	_, ok := contentTypes[s]
	return ok
}

func ExtractMetadataFromContextWithKey(key string, ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(key); len(values) > 0 {
			return values[0], nil
		}
	}

	return "", NewBadRequestError("invalid metadata")
}
