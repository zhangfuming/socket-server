package httpserver

import(
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"../socketserver"
)

type Result struct{
	Ret int
	Reason string
	Data interface{}
}

type ajaxController struct {
}

func StartHttpServer(addr string,flag chan bool){
	http.HandleFunc("/push/message",pushMsg)
	Log("start http server success on ", addr)
	if err := http.ListenAndServe(addr,nil); err != nil{
		log.Fatal("Faile to start http server on ", addr, err)
		flag <- true
	}
}

func pushMsg(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")
	err := r.ParseForm()
	if err != nil {
		outputJson(w, 0, "参数错误", nil)
		return
	}
	message := r.FormValue("message")
	if message == ""{
		outputJson(w, 0, "参数错误", nil)
		return
	}
	socketserver.Broadcast([]byte(message))
	outputJson(w,1,"操作成功",message)
}

func outputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}