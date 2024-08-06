package utils

import (
	"fmt"
	"math/rand"

	"github.com/h2non/filetype"
)

var (
	IMAGE = []string{"jpg", "jpeg", "png", "heic"}
	PDF   = []string{"pdf"}

	SIZE_1KB   = 1 << 10
	SIZE_5KB   = 5 << 10
	SIZE_10KB  = 10 << 10
	SIZE_20KB  = 20 << 10
	SIZE_50KB  = 50 << 10
	SIZE_1MB   = 1 << 20
	SIZE_5MB   = 5 << 20
	SIZE_10MB  = 10 << 20
	SIZE_20MB  = 20 << 20
	SIZE_50MB  = 50 << 20
	SIZE_100MB = 100 << 20
	SIZE_200MB = 200 << 20
	SIZE_500MB = 500 << 20
	SIZE_1GB   = 1 << 30
)

func IsFileType(buf []byte, fileType []string) bool {
	for _, t := range fileType {
		if filetype.Is(buf, t) {
			return true
		}
	}
	return false
}

func FormatSize(size int) string {
	if size < SIZE_1KB {
		return fmt.Sprintf("%d Byte", size)
	} else if size < SIZE_1MB {
		return fmt.Sprintf("%d KB", size/SIZE_1KB)
	} else if size < SIZE_1GB {
		return fmt.Sprintf("%d MB", size/SIZE_1MB)
	} else {
		return fmt.Sprintf("%d GB", size/SIZE_1GB)
	}
}

func GenerateString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
