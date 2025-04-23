package objects

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"storage/infra/dal"
	"storage/internal/apiServer/heartbeat"
	"storage/internal/apiServer/locate"
	rs2 "storage/internal/pkg/rs"
	"storage/myes"
	"time"
)

type File struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}
type SearchResponseData struct {
	Files []File `json:"files"`
}
type SearchResponse struct {
	BaseResp `json:"baseResp"`
	Data     SearchResponseData `json:"data"`
}

// Search es查询
func Search(c *gin.Context) {
	name := c.Query("name")
	fmt.Println("name=", name)
	resp := SearchResponse{
		Data: SearchResponseData{
			Files: make([]File, 0),
		},
		BaseResp: BaseResp{
			Code:    0,
			Message: "success",
		},
	}
	files, err := myes.GetFile(name)
	if err != nil {
		resp.BaseResp.Set(1, err.Error())
		log.Println("err = ", err)
		c.JSON(http.StatusOK, resp)
		return
	}
	for _, file := range files {
		resp.Data.Files = append(resp.Data.Files, File{
			Name: file.Name,
			Hash: file.Hash,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// Get 精确取出对象，按照name取
func Get(c *gin.Context) {
	resp := &BaseResp{}
	hash, err := url.QueryUnescape(c.Query("hash"))
	if err != nil {
		resp.Set(1, "文件名称解析出错")
		c.JSON(http.StatusOK, resp)
		return
	}
	file := dal.Get(hash) //	从MySQL中取出来
	if file.Hash == "" {
		resp.Set(1, "未找到该文件")
		c.JSON(http.StatusOK, resp)
		return
	}
	hash = url.PathEscape(file.Hash) //对字符串转义，保证传输过程中的安全
	stream, e := GetStream(hash, file.Size)
	if e != nil {
		log.Println(e)
		resp.Set(1, e.Error())
		c.JSON(http.StatusOK, resp)
		return
	}
	defer stream.Close()

	//创建临时文件
	f, e := os.CreateTemp("", hash+time.Now().Format(time.DateTime))
	if e != nil {
		resp.Set(1, e.Error())
		c.JSON(http.StatusOK, resp)
		return
	}
	defer os.Remove(f.Name())

	if f != nil {
		if _, err = io.Copy(f, stream); err != nil {
			return
		}
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
		//fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		resp.Set(0, "success")
		c.File(f.Name())
	}
}

func GetStream(hash string, size int64) (*rs2.RSGetStream, error) {
	locateInfo := locate.Locate(hash)
	if len(locateInfo) < rs2.DATA_SHARDS { // 如果定位的数据服务节点数少于数据恢复的最低数量，则无法恢复完整数据，返回定位失败错误
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	// 选择数据服务节点，用于数据恢复
	dataServers := make([]string, 0)
	if len(locateInfo) != rs2.ALL_SHARDS { //从其他数据节点中选择一些节点来进行数据恢复
		dataServers = heartbeat.ChooseRandomDataServers(rs2.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs2.NewRSGetStream(locateInfo, dataServers, hash, size)
}
