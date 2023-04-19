package heartbeat

import (
	"math/rand"
)

// ChooseRandomDataServers
/**
 * @Description: 随机选取一个数据节点服务器返回出去
 * @param n
 * @param exclude
 * @return ds
 */
func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	servers := GetDataServers() //首先我得知道现在有哪些节点存在
	for _, s := range servers {
		_, excluded := reverseExcludeMap[s]
		if !excluded { //获取失败，即该节点不存在于exclude中
			candidates = append(candidates, s)
		}
	}
	length := len(candidates)
	if length < n {
		return
	}
	p := rand.Perm(length)   //将小于length的所有数字随即填充到该数字中
	for i := 0; i < n; i++ { //打乱节点顺序
		ds = append(ds, candidates[p[i]])
	}
	return
}
