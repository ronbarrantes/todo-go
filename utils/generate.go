package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateId() (string, error) {
	var id [6]byte
	_, err := rand.Read(id[:])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", id), nil
}
