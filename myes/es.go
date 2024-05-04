package myes

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/siddontang/go-log/log"
	"reflect"
	"strconv"
)

type EsFile struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
	Size int32  `json:"size"`
}

var (
	Start     = make(chan struct{})
	client    *elastic.Client
	url       = "http://localhost:9200"
	ctx       = context.Background()
	fileIndex = "file"
	mapping   = `{
	  "mappings": {
		"properties": {
		  "name": {
			"type": "text",
			"analyzer": "ik_max_word"
		  },
		  "size": {
			"type": "integer",
			"index": false
		  },
		  "hash": {
			"type": "keyword",
			"index": false
		  }
		}
	  }
	}`
)

func Init() {
	var err error
	client, err = elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(url),
	)
	if err != nil {
		log.Error("连接es失败", err)
		panic(err)
	}
	creatIndex()
	Start <- struct{}{}
}

func creatIndex() {
	exist, err := client.IndexExists(fileIndex).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exist {
		createIndex, err := client.CreateIndex(fileIndex).BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			log.Println("Not acknowledged")
		}
	} else {
		log.Println("already exist")
	}
}

// AddFile 新增文件
func AddFile(size int32, id int64, name, hash string) (err error) {
	file := EsFile{Name: name, Hash: hash, Size: size}
	log.Println("AddFile client = ", client, " ctx = ", ctx)
	post, err := client.Index().Index(fileIndex).BodyJson(file).Id(strconv.Itoa(int(id))).Do(ctx)
	if err != nil {
		log.Println("PostFile err = ", err)
		return err
	}
	log.Println("PostFile success , post = ", post)
	return nil
}

// PutFile 更新文件，暂不使用
func PutFile(id int, name string) (err error) {
	_, err = client.Update().
		Index(fileIndex).
		Id(strconv.Itoa(id)).
		Doc(map[string]string{"name": name}).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// GetFile 按照name作为索引返回一个元数据
func GetFile(name string) (files []*EsFile, err error) {
	termQuery := elastic.NewTermsQuery("name", name)
	result, err := client.Search(fileIndex).Query(termQuery).Do(ctx)
	if err != nil {
		log.Println("Error getting document: ", err)
		return nil, errors.New(fmt.Sprintf("Error getting document: %s ", err))
	}
	if result.TotalHits() == 0 {
		return
	}
	files = make([]*EsFile, result.TotalHits())
	for i, v := range result.Each(reflect.TypeOf(&EsFile{})) {
		files[i] = v.(*EsFile)
	}
	return files, err
}

// DeleteFile MySQL是假删除，但是对于es要真实地删除
func DeleteFile(id int64) error {
	delete, err := client.Delete().Index(fileIndex).Id(strconv.Itoa(int(id))).Do(ctx)
	if err != nil {
		log.Println("DeleteFile err = ", err)
		return err
	}
	log.Println("DeleteFile success, delete = ", delete)
	return nil
}
