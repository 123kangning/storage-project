package rs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"project/go-object-storage/src/lib/objectstream"
	"project/go-object-storage/src/lib/utils"
)

type resumableToken struct {
	Name    string
	Size    int64
	Hash    string
	Servers []string
	Uuids   []string
}

type RSResumablePutStream struct {
	*RSPutStream
	*resumableToken
}

// NewRSResumablePutStream 将元数据 name,size,hash,[]server与[]uuid 放入Token中，将Token和RSPutStream集成到RSResumablePutStream中
func NewRSResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, e := NewRSPutStream(dataServers, hash, size)
	if e != nil {
		return nil, e
	}
	uuids := make([]string, ALL_SHARDS)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
	}
	token := &resumableToken{name, size, hash, dataServers, uuids}
	return &RSResumablePutStream{putStream, token}, nil
}

// NewRSResumablePutStreamFromToken 用Token还原出对应的的文件流，写入文件时同样是向dataServer发送patch请求
func NewRSResumablePutStreamFromToken(token string) (*RSResumablePutStream, error) {
	b, e := base64.StdEncoding.DecodeString(token)
	if e != nil {
		return nil, e
	}

	var t resumableToken
	e = json.Unmarshal(b, &t)
	if e != nil {
		return nil, e
	}

	writers := make([]io.Writer, ALL_SHARDS)
	for i := range writers {
		writers[i] = &objectstream.TempPutStream{Server: t.Servers[i], Uuid: t.Uuids[i]}
	}
	enc := NewEncoder(writers)
	return &RSResumablePutStream{&RSPutStream{enc}, &t}, nil
}

// ToToken 之后用Token进行分片发送文件
func (s *RSResumablePutStream) ToToken() string {
	b, _ := json.Marshal(s)
	return base64.StdEncoding.EncodeToString(b)
}

func (s *RSResumablePutStream) CurrentSize() int64 {
	r, e := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if e != nil {
		log.Println(e)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.StatusCode)
		return -1
	}
	size := utils.GetSizeFromHeader(r.Header) * DATA_SHARDS
	if size > s.Size {
		size = s.Size
	}
	return size
}
