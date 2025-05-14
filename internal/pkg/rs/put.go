package rs

import (
	"fmt"
	"io"
	"storage/internal/pkg/objectstream"
)

type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServers []string, hash string, size int) (*RSPutStream, error) {
	if len(dataServers) != ALL_SHARDS { // 如果参数中的数据服务节点数量小于最低标准，则无法存储数据
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS //数据分片，每一个数据片的大小
	writers := make([]io.Writer, ALL_SHARDS)
	var e error
	for i := range writers {
		writers[i], e = objectstream.NewTempPutStream(dataServers[i],
			fmt.Sprintf("%s.%d", hash, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	//生成RS编码器对象
	enc := NewEncoder(writers)
	return &RSPutStream{enc}, nil
}

func (s *RSPutStream) Commit(success bool) {
	s.Flush()
	for i := range s.writers {
		//成功时，向dataServer发送put请求，转正文件数据
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}
