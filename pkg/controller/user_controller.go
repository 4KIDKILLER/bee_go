package controller

import "net/http"

func UserController(mux *http.ServeMux) {
	mux.HandleFunc("GET /getUserInfo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("获取用户信息"))
	})
}
