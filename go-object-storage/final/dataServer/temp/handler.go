package temp

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	// PUT 将临时文件转换成正式文件
	if m == http.MethodPut {
		put(w, r)
		return
	}
	// PATCH 写入数据到临时文件
	if m == http.MethodPatch {
		patch(w, r)
		return
	}
	// POST 创建临时文件
	if m == http.MethodPost {
		post(w, r)
		return
	}
	// DELETE 删除临时文件
	if m == http.MethodDelete {
		del(w, r)
		return
	}
	// HEAD 查看文件状态
	if m == http.MethodHead {
		head(w, r)
		return
	}
	// GET 查看文件状态
	if m == http.MethodGet {
		get(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
