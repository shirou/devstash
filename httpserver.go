package main

import (
	"net/http"
	"fmt"
	"html"
	"strconv"
	"html/template"
)

type DevStashHTTP struct {
	Config Config
}


func (d DevStashHTTP) storeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func (d DevStashHTTP) listHandler(w http.ResponseWriter, r *http.Request) {
	path := d.Config.Server.Directory
	max_results := 100  // TODO: temporary.
	finfo_list, err := readIndexFile(path, max_results)
	if err != nil{
		http.Error(w, "Could not read Index File", 500)
	}

	t, _ := template.ParseFiles("list.html")
	t.Execute(w, finfo_list)
}

func (d DevStashHTTP) showHandler(w http.ResponseWriter, r *http.Request) {
}


func (d DevStashHTTP) StartHTTPServer() {
	port := d.Config.Server.Port
//	store_dir := config.Server.Directory

	http.HandleFunc("/p", d.storeHandler)
	http.HandleFunc("/index.html", d.listHandler)
	http.HandleFunc("/", d.showHandler)

	addr := ":" + strconv.Itoa(port)
	fmt.Println("Start at " + addr)
	if d.Config.Server.Ssl == true{
		fmt.Println(http.ListenAndServeTLS(addr, d.Config.Server.CertFile, d.Config.Server.KeyFile, nil))
	}else{
		fmt.Println(http.ListenAndServe(addr, nil))
	}
}

