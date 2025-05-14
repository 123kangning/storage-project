package errs

// 1xx 文件类
// 2xx 用户类

var (
	ErrFileShardLossTooMany = NewErrorCode(1001, "文件分片丢失过多，下载失败")
)
