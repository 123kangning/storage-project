package objects

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"storage/infra/dal"
	"storage/infra/myes"
	"storage/internal/apiServer/heartbeat"
	"storage/internal/apiServer/locate"
	"storage/internal/pkg/errs"
	rs2 "storage/internal/pkg/rs"
	"strconv"
	"time"
)

type File struct {
	Name   string `json:"name"`
	Hash   string `json:"hash"`
	Size   int    `json:"size"`
	Source string `json:"source,omitempty"`
}
type SearchResponseData struct {
	Files []File `json:"files"`
	Total int64  `json:"total"`
}
type SearchResponse struct {
	BaseResp `json:"baseResp"`
	Data     SearchResponseData `json:"data"`
}

// GetUserUploadedFiles 获取用户自己上传的所有文件
func GetUserUploadedFiles(c *gin.Context) {
	resp := &SearchResponse{
		Data: SearchResponseData{
			Files: make([]File, 0),
		},
		BaseResp: BaseResp{
			Code:    0,
			Message: "success",
		},
	}
	name := c.Query("name")
	// 获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		resp.BaseResp.Set(1, "未获取到用户 ID")
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		resp.BaseResp.Set(1, "用户 ID 类型错误")
		c.JSON(http.StatusUnauthorized, resp)
		return
	}
	from, size := page(c)
	files, total, err := myes.GetFile(name, from, size, userIDInt64)
	resp.Data.Total = total
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
			Size: file.Size,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// Search es查询
func Search(c *gin.Context) {
	name := c.Query("name")
	from, size := page(c)

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
	files, total, err := myes.GetFile(name, from, size, 0)
	resp.Data.Total = total
	if err != nil {
		resp.BaseResp.Set(1, err.Error())
		log.Println("err = ", err)
		c.JSON(http.StatusOK, resp)
		return
	}
	userIdList := make([]int64, 0)
	for _, file := range files {
		userIdList = append(userIdList, file.Source)
	}
	userId2Username, err := dal.GetUserIDToUsernameMap(userIdList)
	if err != nil {
		resp.BaseResp.Set(1, err.Error())
		log.Println("err = ", err)
		c.JSON(http.StatusOK, resp)
		return
	}
	for _, file := range files {
		resp.Data.Files = append(resp.Data.Files, File{
			Name:   file.Name,
			Hash:   file.Hash,
			Size:   file.Size,
			Source: userId2Username[file.Source],
		})
	}
	c.JSON(http.StatusOK, resp)
}

func page(c *gin.Context) (int, int) {
	// 获取分页参数，默认从第 1 条开始，每页 10 条记录
	fromStr := c.DefaultQuery("from", "1")
	sizeStr := c.DefaultQuery("size", "10")
	from, err := strconv.Atoi(fromStr)
	if err != nil {
		from = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		size = 10
	}
	from = (from - 1) * size
	return from, size
}

// Down 精确取出对象，按照name取
func Down(c *gin.Context) {
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
		if errors.Is(e, errs.ErrFileShardLossTooMany) {
			err = dal.Delete(hash, 0)
			if err != nil {
				resp.Set(1, err.Error())
				c.JSON(http.StatusOK, resp)
				return
			}
		}
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

func GetStream(hash string, size int) (*rs2.RSGetStream, error) {
	locateInfo := locate.Locate(hash)
	if len(locateInfo) < rs2.DATA_SHARDS { // 如果定位的数据服务节点数少于数据恢复的最低数量，则无法恢复完整数据，返回定位失败错误
		log.Println("object", hash, "locate fail , locateInfo=", locateInfo)
		return nil, errs.ErrFileShardLossTooMany
	}
	// 选择数据服务节点，用于数据恢复
	dataServers := make([]string, 0)
	if len(locateInfo) != rs2.ALL_SHARDS { //从其他数据节点中选择一些节点来进行数据恢复
		dataServers = heartbeat.ChooseRandomDataServers(rs2.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs2.NewRSGetStream(locateInfo, dataServers, hash, size)
}
