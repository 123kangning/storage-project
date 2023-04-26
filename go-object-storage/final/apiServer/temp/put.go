package temp

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"project/go-object-storage/final/apiServer/locate"
	"project/go-object-storage/src/lib/es"
	"project/go-object-storage/src/lib/rs"
	"project/go-object-storage/src/lib/utils"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset := utils.GetOffsetFromHeader(r.Header)
	if current != offset {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	bytes := make([]byte, rs.BLOCK_SIZE)
	for {
		n, e := io.ReadFull(r.Body, bytes)
		if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current += int64(n)
		if current > stream.Size {
			stream.Commit(false)
			log.Println("resumable put exceed size")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if n != rs.BLOCK_SIZE && current != stream.Size {
			return
		}
		//写入enc.cache缓冲区中
		stream.Write(bytes[:n])
		if current == stream.Size {
			//发送patch调用，写入dataServer
			stream.Flush()
			getStream, e := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			hash := url.PathEscape(utils.CalculateHash(getStream))
			if hash != stream.Hash {
				//回退
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if locate.Exist(url.PathEscape(hash)) { //判断至少四个数据节点是否保存了该数据
				stream.Commit(false) //回退，已经保存了该数据
			} else {
				//提交
				stream.Commit(true)
			}
			e = es.AddVersion(stream.Name, stream.Hash, stream.Size)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}
