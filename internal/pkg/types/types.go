package types

// LocateMessage dataServer的地址和其存储的文件分片id,通过hash可以索引到这个结构体信息
type LocateMessage struct {
	Addr string
	Id   int
}
