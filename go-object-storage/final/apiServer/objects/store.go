package objects

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"project/go-object-storage/final/apiServer/locate"
	"project/go-object-storage/src/lib/utils"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}

	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}
