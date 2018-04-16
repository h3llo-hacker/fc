package utils

import (
	"encoding/hex"

	"github.com/h3llo-hacker/fc/config"

	"golang.org/x/crypto/scrypt"
)

func Password(oriPass string) string {
	salt := config.Salt
	dk, _ := scrypt.Key([]byte(oriPass), []byte(salt), 16384, 8, 1, 32)
	encPass := hex.EncodeToString(dk)
	return encPass
}
