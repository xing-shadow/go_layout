package utils

import (
	"github.com/google/uuid"
	"strings"
)

func GetUUid() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
