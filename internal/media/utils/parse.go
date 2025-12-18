package utils

import (
	"strconv"

	"github.com/rs/zerolog/log"
)

func ParseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Error().Msgf("failed to parse to bool: %s", s)
	}
	return b
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Error().Msgf("failed to parse to bool: %s", s)
	}
	return i
}
