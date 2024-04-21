package rs

import (
	"fmt"
	"io"
	objectstream2 "storage/internal/pkg/objectstream"
)

type RSGetStream struct {
	*decoder
}

// NewRSGetStream 参数 locateInfo 为读取切片的数据服务节点，dataServers 为恢复切片的数据服务节点 ，hash 为数据散列值，size 为数据总长度
func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	readers := make([]io.Reader, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		server := locateInfo[i]
		if server == "" { // 若节点不存在则取随机节点补充
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		//真正调用data层的API,获取数据,若节点存在则读取该节点保存的切片数据
		reader, e := objectstream2.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if e == nil {
			readers[i] = reader
		}
	}
	//恢复分片，当节点不存在时取随机节点补充，或者是节点存在时读取切片数据抛出异常都会导致 readers 某个值为 nil
	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	for i := range readers {
		if readers[i] == nil { // 当切片元素为 nil 时，创建临时对象写入流用于恢复分片(RS还原出文件之后就可以通过writer流调用patch将数据写回dataServer中)
			writers[i], e = objectstream2.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if e != nil {
				return nil, e
			}
		}
	}

	dec := NewDecoder(readers, writers, size)
	return &RSGetStream{dec}, nil
}

// Close 将临时对象转正
func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream2.TempPutStream).Commit(true)
		}
	}
}

func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")
	}
	if offset < 0 {
		panic("only support forward seek")
	}
	for offset != 0 {
		length := int64(BLOCK_SIZE)
		if offset < length {
			length = offset
		}
		buf := make([]byte, length)
		io.ReadFull(s, buf)
		offset -= length
	}
	return offset, nil
}
