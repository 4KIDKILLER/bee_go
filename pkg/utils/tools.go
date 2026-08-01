package utils

import (
	"fmt"
	"strings"
	"time"

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

func GetUploadPath() string {
	now := time.Now()
	return fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day())
}
