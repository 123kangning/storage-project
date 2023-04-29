package es

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"runtime"
)

type Metadata struct { //元数据结构体 包含名称、版本、size、哈希值
	Name    string
	Version int
	Size    int64
	Hash    string
}

var (
	client *elastic.Client
	url    = "http://localhost:9200"
	ctx    = context.Background()
)

func Init() {
	var err error
	client, err = elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Println("连接es失败", err)
	}
	log.Println(elastic.Version)
	log.Println(client)
}

// 按照名称+版本作为索引返回一个元数据
func getMetadata(name string, version int) (meta Metadata, e error) {
	index := fmt.Sprintf("%s_%d", name, version)
	indexExists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Println("Error checking index existence: ", err)
		return Metadata{}, errors.New(fmt.Sprintf("Error checking index existence: %s ", err))
	} else if !indexExists {
		log.Println("Index doesn't exist: ", index)
		return Metadata{}, errors.New(fmt.Sprintf("Index doesn't exist: %s ", index))
	}
	result, err := client.Get().Index(index).Id("1").Do(ctx)
	if err != nil {
		log.Println("Error getting document: ", err)
		return Metadata{}, errors.New(fmt.Sprintf("Error getting document: %s ", err))
	}
	b := result.Source
	err = json.Unmarshal(b, &meta)
	if err != nil {
		log.Println("Unmarshal error ", err)
	}
	return meta, err
}

func SearchLatestVersion(name string) (meta Metadata, err error) {
	defer func() {
		if err := recover(); err != nil {
			// 获取 panic 的堆栈信息
			var stackTrace string
			for i := 1; ; i++ {
				pc, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				stackTrace += fmt.Sprintf("%s:%d (0x%x)\n", file, line, pc)
			}

			// 获取 panic 的错误信息
			errMsg := fmt.Sprintf("%v", err)

			// 打印堆栈信息和错误信息
			fmt.Printf("panic: %s\n%s", errMsg, stackTrace)
		}
	}()
	searchService := client.Search()
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("name", name))
	searchSource.Sort("version", false) // 按照 version 字段降序排序
	searchSource.Size(1)
	searchResult, err := searchService.SearchSource(searchSource).Do(ctx)
	if err != nil {
		log.Println("SearchLatestVersion ", err)
	}
	fmt.Println("searchResult.Hits = ", searchResult.Hits)
	//没找到，也就是第一次查看时
	if len(searchResult.Hits.Hits) == 0 {
		return Metadata{}, nil
	}
	b := searchResult.Hits.Hits[0].Source
	err = json.Unmarshal(b, &meta)
	if err != nil {
		log.Println("Unmarshal error ", err)
	}
	return
}

// GetMetadata 对getMetadata的一层封装，增加了未指定版本时候自动获取最新的版本功能
func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

// PutMetadata 把元数据存进去，信息包括名称，版本，size和哈希值
func PutMetadata(name string, version int, size int64, hash string) error {
	index := fmt.Sprintf("%s_%d", name, version)
	doc := map[string]interface{}{
		"name":    name,
		"version": version,
		"size":    size,
		"hash":    hash,
	}
	_, err := client.Index().Index(index).Id("1").BodyJson(doc).Do(ctx)
	if err != nil {
		fmt.Println("Error indexing document: ", err)
	} else {
		fmt.Println("Document indexed:", doc)
	}
	log.Println("put success")
	return err
}

// AddVersion 版本迭代 version+1
func AddVersion(name, hash string, size int64) error {
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

// SearchAllVersions 找出这个对象存在的所有版本
func SearchAllVersions(name string, from, size int) (metas []Metadata, err error) {
	searchService := client.Search()
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("name", name))
	searchSource.Sort("version", false) // 按照 version 字段降序排序
	searchSource.From(from)
	searchSource.Size(size)
	searchResult, err := searchService.SearchSource(searchSource).Do(ctx)
	if err != nil {
		log.Println("SearchLatestVersion ", err)
	}
	fmt.Println("searchResult.Hits.Hits = ", searchResult.Hits.Hits)
	metas = make([]Metadata, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		b := hit.Source
		_ = json.Unmarshal(b, &metas[i])
	}
	return
}
