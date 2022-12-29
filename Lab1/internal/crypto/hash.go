package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

type HashInterface interface {
	GenerateHash(string) string
}

type SHA256Hash struct{}

func (sha *SHA256Hash) GenerateHash(text string) string {
	byteHash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(byteHash[:])
}

type SHA512Hash struct{}

func (sha *SHA512Hash) GenerateHash(text string) string {
	byteHash := sha512.Sum512([]byte(text))
	return hex.EncodeToString(byteHash[:])
}

type MD5Hash struct{}

func (md *MD5Hash) GenerateHash(text string) string {
	byteHash := md5.Sum([]byte(text))
	return hex.EncodeToString(byteHash[:])
}
