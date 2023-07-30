package myes

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"reflect"
	"strconv"
)

//type Metadata struct { //元数据结构体 包含名称、版本、size、哈希值
//	Name    string
//	Version int
//	Size    int64
//	Hash    string
//}

type File struct {
	name   string
	author string
	hash   string
	size   int
}

var (
	client    *elastic.Client
	url       = "http://localhost:9200"
	ctx       = context.Background()
	fileIndex = "file"
	mapping   = `{
  "mappings": {
    "properties": {
      "name":{
        "type": "text",
        "copy_to": "all"
      },
      "author":{
        "type": "keyword",
        "copy_to": "all"
      },
      "hash":{
        "type": "keyword",
        "index": false
      },
      "size":{
        "type": "integer",
        "index": false
      },
      "all":{
        "type": "text",
        "analyzer": "ik_max_word"
      }
    }
  }
}`
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

// PostFile 新增文件
func PostFile(size, id int, name, hash string) (err error) {
	file := File{name: name, hash: hash, size: size}
	post, err := client.Index().Index(fileIndex).BodyJson(file).Id(strconv.Itoa(id)).Do(ctx)
	if err != nil {
		log.Println("PostFile err = ", err)
		return err
	}
	log.Println("PostFile success , post = ", post)
	return nil
}

// PutFile 更新文件
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

// GetAllFile 获取所有文件
func GetAllFile() (files []*File, err error) {
	query := elastic.NewMatchAllQuery()
	searchResult, err := client.Search().
		Index(fileIndex).
		Query(query).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if searchResult.TotalHits() == 0 {
		return nil, errors.New("Search 出错 ")
	}
	files = make([]*File, searchResult.TotalHits())
	for i, v := range searchResult.Each(reflect.TypeOf(&File{})) {
		files[i] = v.(*File)
	}
	return files, nil
}

// GetFile 按照name或author作为索引返回一个元数据
func GetFile(name string) (files []*File, err error) {
	termQuery := elastic.NewTermsQuery("all", name)
	result, err := client.Search(fileIndex).Query(termQuery).Do(ctx)
	if err != nil {
		log.Println("Error getting document: ", err)
		return nil, errors.New(fmt.Sprintf("Error getting document: %s ", err))
	}
	if result.TotalHits() == 0 {
		return
	}
	files = make([]*File, result.TotalHits())
	for i, v := range result.Each(reflect.TypeOf(&File{})) {
		files[i] = v.(*File)
	}
	return files, err
}

//func SearchLatestVersion(name string) (meta Metadata, err error) {
//	defer func() {
//		if err := recover(); err != nil {
//			// 获取 panic 的堆栈信息
//			var stackTrace string
//			for i := 1; ; i++ {
//				pc, file, line, ok := runtime.Caller(i)
//				if !ok {
//					break
//				}
//				stackTrace += fmt.Sprintf("%s:%d (0x%x)\n", file, line, pc)
//			}
//
//			// 获取 panic 的错误信息
//			errMsg := fmt.Sprintf("%v", err)
//
//			// 打印堆栈信息和错误信息
//			fmt.Printf("panic: %s\n%s", errMsg, stackTrace)
//		}
//	}()
//	searchService := client.Search()
//	searchSource := elastic.NewSearchSource()
//	searchSource.Query(elastic.NewMatchQuery("name", name))
//	searchSource.Sort("version", false) // 按照 version 字段降序排序
//	searchSource.Size(1)
//	searchResult, err := searchService.SearchSource(searchSource).Do(ctx)
//	if err != nil {
//		log.Println("SearchLatestVersion ", err)
//	}
//	fmt.Println("searchResult.Hits = ", searchResult.Hits)
//	//没找到，也就是第一次查看时
//	if len(searchResult.Hits.Hits) == 0 {
//		return Metadata{}, nil
//	}
//	b := searchResult.Hits.Hits[0].Source
//	err = json.Unmarshal(b, &meta)
//	if err != nil {
//		log.Println("Unmarshal error ", err)
//	}
//	return
//}

//// GetMetadata 对getMetadata的一层封装，增加了未指定版本时候自动获取最新的版本功能
//func GetMetadata(name string, version int) (Metadata, error) {
//	if version == 0 {
//		return SearchLatestVersion(name)
//	}
//	return getMetadata(name, version)
//}

//// PutMetadata 把元数据存进去，信息包括名称，版本，size和哈希值
//func PutMetadata(name string, version int, size int64, hash string) error {
//	index := fmt.Sprintf("%s_%d", name, version)
//	doc := map[string]interface{}{
//		"name":    name,
//		"version": version,
//		"size":    size,
//		"hash":    hash,
//	}
//	_, err := client.Index().Index(index).Id("1").BodyJson(doc).Do(ctx)
//	if err != nil {
//		fmt.Println("Error indexing document: ", err)
//	} else {
//		fmt.Println("Document indexed:", doc)
//	}
//	log.Println("put success")
//	return err
//}

//// AddVersion 版本迭代 version+1
//func AddVersion(name, hash string, size int64) error {
//	version, e := SearchLatestVersion(name)
//	if e != nil {
//		return e
//	}
//	return PutMetadata(name, version.Version+1, size, hash)
//}
//
//// SearchAllVersions 找出这个对象存在的所有版本
//func SearchAllVersions(name string, from, size int) (metas []Metadata, err error) {
//	searchService := client.Search()
//	searchSource := elastic.NewSearchSource()
//	searchSource.Query(elastic.NewMatchQuery("name", name))
//	searchSource.Sort("version", false) // 按照 version 字段降序排序
//	searchSource.From(from)
//	searchSource.Size(size)
//	searchResult, err := searchService.SearchSource(searchSource).Do(ctx)
//	if err != nil {
//		log.Println("SearchLatestVersion ", err)
//	}
//	fmt.Println("searchResult.Hits.Hits = ", searchResult.Hits.Hits)
//	metas = make([]Metadata, len(searchResult.Hits.Hits))
//	for i, hit := range searchResult.Hits.Hits {
//		b := hit.Source
//		_ = json.Unmarshal(b, &metas[i])
//	}
//	return
//}
