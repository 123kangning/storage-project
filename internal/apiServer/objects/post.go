package objects

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"storage/infra/dal"
	"storage/internal/apiServer/heartbeat"
	rs2 "storage/internal/pkg/rs"
	"storage/internal/pkg/utils"
)

type BaseResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (r *BaseResp) Set(code int, message string) {
	r.Code = code
	r.Message = message
}
func Post(c *gin.Context) {
	log.Println("api.object.put")
	resp := &BaseResp{}

	// 获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		resp.Set(1, "未获取到用户 ID")
		c.JSON(http.StatusUnauthorized, resp)
		return
	}
	userIDInt64, ok := userID.(int64)
	if !ok {
		resp.Set(1, "用户 ID 类型错误")
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Println("file error , ", err)
		resp.Set(1, fmt.Sprintln("file error , ", err))
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	log.Println("from file complete, name =", file.Filename, ",len =", len(file.Filename))

	r, err := file.Open()

	if err != nil {
		log.Println("file error , ", err)
		resp.Set(1, fmt.Sprintln("file error , ", err))
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	defer r.Close()
	//get hash
	hash := utils.CalculateHash(r)

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		log.Println("file Seek error , ", err)
		resp.Set(1, fmt.Sprintln("file Seek error , ", err))
		return
	} //重置reader的位置

	resFile := dal.Get(hash)
	if resFile.Hash == hash {
		resp.Set(1, "文件已存在")
		c.JSON(http.StatusOK, resp)
		return
	}

	code, e := storeObject(r, hash, int(file.Size)) //存储文件到/objects,返回状态码以及error
	if e != nil || code != http.StatusOK {
		log.Println(e)
		resp.Set(1, e.Error())
		c.JSON(code, resp)
		return
	}

	name := file.Filename                                //组成名字
	e = dal.Add(name, hash, int(file.Size), userIDInt64) //更新数据库
	if e != nil {
		log.Println(e)
		resp.Set(1, e.Error())
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Set(0, "success")
	log.Println("resp = ", resp)
	c.JSON(http.StatusOK, resp)
}

// 调用dataServer生成文件，暂时还未写入，返回writer
func putStream(hash string, size int) (*rs2.RSPutStream, error) {
	log.Println("api.objects.putStream")
	// 获取全部的数据服务节点，无需排除任何节点
	servers := heartbeat.ChooseRandomDataServers(rs2.ALL_SHARDS, nil)
	// 如果数据服务节点长度不等于分片长度，则无法完整保存数据，提示报错
	if len(servers) != rs2.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}
	return rs2.NewRSPutStream(servers, hash, size) //调用dataServer生成文件，暂时还未写入
}
