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

// GetFile 按照name作为索引返回分页后的元数据
func GetFile(name string, from int, size int) (files []*EsFile, total int64, err error) {
	// 构建查询条件
	termQuery := elastic.NewTermsQuery("name", name)
	// 执行 count 查询
	fmt.Println("client=", client, "ctx=", ctx)
	count, err := client.Count(fileIndex).Query(termQuery).Do(ctx)
	if err != nil {
		log.Println("Error counting documents: ", err)
		return nil, 0, errors.New(fmt.Sprintf("Error counting documents: %s ", err))
	}
	total = count

	// 执行搜索查询
	result, err := client.Search(fileIndex).
		Query(termQuery).
		From(from). // 设置起始位置
		Size(size). // 设置每页大小
		Do(ctx)
	if err != nil {
		log.Println("Error getting document: ", err)
		return nil, 0, errors.New(fmt.Sprintf("Error getting document: %s ", err))
	}
	if result.TotalHits() == 0 {
		return
	}
	files = make([]*EsFile, 0, result.Hits.TotalHits.Value)
	for _, v := range result.Each(reflect.TypeOf(&EsFile{})) {
		if file, ok := v.(*EsFile); ok {
			files = append(files, file)
		}
	}
	return files, total, nil
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
