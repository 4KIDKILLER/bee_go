package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GetUUID() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	id := strings.ReplaceAll(uid.String(), "-", "")

	return id, nil
}
