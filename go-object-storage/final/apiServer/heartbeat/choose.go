package heartbeat

import (
	"log"
	"math/rand"
)

// ChooseRandomDataServers
/*随机选择多个数据服务节点返回，参数 n 指定返回多少个节点，参数 exclude 指定排除哪些节点*/
func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	servers := GetDataServers() //获取全部节点,首先我得知道现在有哪些节点存在
	log.Println("dataServer = ", servers)
	for _, s := range servers {
		_, excluded := reverseExcludeMap[s]
		if !excluded { //只需要获取除exclude之外的节点
			candidates = append(candidates, s)
		}
	}
	length := len(candidates)
	//log.Println("candidates = ",candidates)
	log.Println(" length = ", len(candidates), " n = ", n)
	if length < n {
		return
	}
	p := rand.Perm(length)   //将小于length的所有数字随即填充到该数字中
	for i := 0; i < n; i++ { //打乱节点顺序
		ds = append(ds, candidates[p[i]])
	}
	return
}
