package rs

import (
	"errors"
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
)

type decoder struct {
	readers   []io.Reader
	writers   []io.Writer
	enc       reedsolomon.Encoder
	size      int
	cache     []byte
	cacheSize int
	total     int
}

/* 构建 RS 码解码器 */
func NewDecoder(readers []io.Reader, writers []io.Writer, size int) *decoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &decoder{readers, writers, enc, size, nil, 0, 0}
}

/* 实现 Read 方法 */
func (d *decoder) Read(p []byte) (n int, err error) {
	if d.cacheSize == 0 {
		e := d.getData()
		if e != nil {
			return 0, e
		}
	}
	length := len(p)
	if d.cacheSize < length {
		length = d.cacheSize
	}
	d.cacheSize -= length
	copy(p, d.cache[:length])
	d.cache = d.cache[length:]
	return length, nil
}

/*获取数据*/
func (d *decoder) getData() error {
	// 如果当前已解码的数据等于原始数据大小，则所有数据已经被读取，返回文件尾标识 io.EOF
	if d.total == d.size {
		return io.EOF
	}

	// 创建 []byte 类型的切片，长度为 6 ，用于保存相应切片的数据
	shards := make([][]byte, ALL_SHARDS)
	// 创建一个整型切片，用于保存修复切片的下标
	repairIds := make([]int, 0)

	// 遍历 readers
	for i := range shards {
		// 如果 readers[i] 为空则说明分片数据丢失需要修复
		if d.readers[i] == nil {
			fmt.Println("readers[", i, "] == nil")
			repairIds = append(repairIds, i)
		} else {
			// 如果 readers[i] 不为空说明分片数据正常，则将分片数据保存到 shards[i] 中
			shards[i] = make([]byte, BLOCK_PER_SHARD)
			/*
				ReadFull准确地将len(buf)字节从r读入buf。它返回复制的字节数，如果读取的字节数较少，则返回一个错误。只有在没有读取字节的情况下，错误才为EOF。
				如果在读取部分但不是所有字节后发生EOF, ReadFull将返回ErrUnexpectedEOF。
				返回时，n == len(buf)当且仅当err == nil时。
				如果r返回至少读取len(buf)字节的错误，则删除该错误。
			*/
			n, e := io.ReadFull(d.readers[i], shards[i])
			if e != nil && e != io.EOF && !errors.Is(e, io.ErrUnexpectedEOF) {
				shards[i] = nil
			} else if n != BLOCK_PER_SHARD {
				shards[i] = shards[i][:n]
			}
		}
	}
	// 尝试重构已丢失的数据切片
	/*
		如果可能，重建将重新创建丢失的碎片。给定一个碎片列表，其中一些包含数据，填充那些没有数据的碎片。
		数组的长度必须等于shard的总数。您可以通过将其设置为nil或零长度来指示一个碎片丢失。
		如果shard长度为零，但有足够的容量，该内存将被使用，否则将分配一个新的[]字节。
		如果碎片太少，无法重建丢失的碎片，则返回ErrTooFewShards。重构的碎片集是完整的，但完整性未验证。使用Verify函数检查数据集是否正确。
	*/
	e := d.enc.Reconstruct(shards)
	if e != nil {
		return e
	}
	for i := range repairIds {
		id := repairIds[i]
		d.writers[id].Write(shards[id])
	}
	for i := 0; i < DATA_SHARDS; i++ {
		shardSize := len(shards[i])
		if d.total+shardSize > d.size {
			shardSize -= d.total + shardSize - d.size
		}
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.cacheSize += shardSize
		d.total += shardSize
	}
	return nil
}
