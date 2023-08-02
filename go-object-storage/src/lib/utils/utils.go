package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// CalculateHash 计算该文件内容的hash值（文件内容从r中读取）
func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return hex.EncodeToString(h.Sum(nil))
}
