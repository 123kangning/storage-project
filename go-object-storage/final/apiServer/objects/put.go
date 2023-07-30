package objects

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"storage/final/apiServer/heartbeat"
	"storage/src/lib/rs"
)

type BaseResp struct {
	statusCode    int
	statusMessage string
}

func (r BaseResp) Set(code int, message string) {
	r.statusCode = code
	r.statusMessage = message
}
func Put(c *gin.Context) {
	log.Println("api.object.put")
	resp := &BaseResp{}

	file, err := c.FormFile("file")
	r, err := file.Open()

	if err != nil {
		log.Println("file error , ", err)
		resp.Set(1, fmt.Sprintln("file error , ", err))
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	//get hash
	hash := c.GetHeader("hash")
	code, e := storeObject(r, hash, file.Size) //存储文件到/objects,返回状态码以及error
	if e != nil || code != http.StatusOK {
		log.Println(e)
		resp.Set(1, e.Error())
		c.JSON(code, resp)
		return
	}

	name := file.Filename                      //组成名字
	e = myes.AddVersion(name, hash, file.Size) //更新版本
	if e != nil {
		log.Println(e)
		resp.Set(1, e.Error())
		c.JSON(http.StatusInternalServerError, resp)
	}
}

// 调用dataServer生成文件，暂时还未写入，返回writer
func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	log.Println("api.objects.putStream")
	// 获取全部的数据服务节点，无需排除任何节点
	servers := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	// 如果数据服务节点长度不等于分片长度，则无法完整保存数据，提示报错
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}
	return rs.NewRSPutStream(servers, hash, size) //调用dataServer生成文件，暂时还未写入
}
