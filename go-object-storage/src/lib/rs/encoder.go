package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

/* RS 码编码器结构体 */
type encoder struct {
	writers []io.Writer //这个Writer接口类型，由TempPutStream类型实现
	enc     reedsolomon.Encoder
	cache   []byte
}

// NewEncoder /* 生成 RS 码编码器 */
func NewEncoder(writers []io.Writer) *encoder {
	// 这里生成了 DATA_SHARDS 个数据分片和 PARITY_SHARDS 个校验分片的 RS 码编码器 enc
	/*
		New创建一个新的编码器，并将其初始化为您想要使用的数据碎片和奇偶碎片的数量。
		您可以重用这个编码器。注意，总的shard的最大数量是256。如果没有提供选项，则使用默认选项。
	*/
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &encoder{writers, enc, nil}
}

//写入encoder.cache缓存，缓存写满之后flush，真正写入dataServer
func (e *encoder) Write(p []byte) (n int, err error) {
	// 获取待写入的数据 p 的总长度
	length := len(p)
	// 当前缓存的数据长度
	current := 0
	// 将待写入的数据 p 以块的形式保存在缓存
	for length != 0 {
		// 块数据长度 - 已缓存数据长度 = 可缓存数据长度
		next := BLOCK_SIZE - len(e.cache)
		// 如果可缓存数据长度比数据长度 p 大，则可以继续缓存数据
		if next > length {
			next = length
		}
		// 新增缓存数据
		e.cache = append(e.cache, p[current:current+next]...)
		if len(e.cache) == BLOCK_SIZE {
			e.Flush()
		}
		// 当前已缓存数据长度增加
		current += next
		// 需要保存数据 p 的长度减少
		length -= next
	}
	return len(p), nil
}

// Flush 将缓存中的数据写入真正的数据切片
func (e *encoder) Flush() {
	// 如果缓存数据长度为 0 则直接返回
	if len(e.cache) == 0 {
		return
	}

	// 调用 Split 方法将缓存数据切分数据片，并且生成两个空的奇偶校验片
	/*
		将一个数据切片分割成给定给编码器的分片数，并创建空的奇偶校验分片。
		数据将被分割成同等大小的碎片。如果数据大小不能被分片数整除，那么最后一个分片将包含额外的零。
		必须至少有一个字节，否则将返回ErrShortData。
		除了最后一个分片外，数据不会被复制，所以以后不应该修改输入切片的数据。
	*/
	shards, _ := e.enc.Split(e.cache)

	// 调用 Encode 方法生成数据分片的奇偶校验片
	/*
		为一组数据碎片编码奇偶性。输入是'shards'，包含数据碎片，后跟奇偶碎片。
		碎片的数量必须匹配给New()的数量。每个shard是一个字节数组，它们必须都是相同的大小。
		奇偶校验碎片将总是被覆盖，而数据碎片将保持不变，所以当它运行时，从数据碎片读取数据是安全的。
	*/
	e.enc.Encode(shards)
	for i := range shards {
		// 每个数据分片文件真正写入分片数据
		e.writers[i].Write(shards[i])
	}
	// 重置缓存大小
	e.cache = []byte{}
}
