package objects

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"storage/internal/dataServer/locate"
	"strings"
)

func getFile(name string) string {
	// 打开磁盘中的文件，检查该文件哈希值是否与请求的哈希值一致
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	if len(files) != 1 {
		return ""
	}
	file := files[0]
	// 计算哈希值并进行 url 转义
	h := sha256.New()
	SendFile(h, file)
	d := url.PathEscape(hex.EncodeToString(h.Sum(nil)))
	hash := strings.Split(file, ".")[2]
	// 如果磁盘文件计算的哈希值与请求的哈希值不一致，说明磁盘的文件数据被降解，需要删除该错误的数据
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		// 从全局缓存中移除该文件
		locate.Del(hash)
		// 磁盘删除该文件
		os.Remove(file)
		return ""
	}
	// 哈希值一致则返回文件路径
	return file
}
