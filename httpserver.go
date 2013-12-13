package main

import (
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type DevStashHTTP struct {
	Config Config
}

func (d DevStashHTTP) storeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func (d DevStashHTTP) listHandler(w http.ResponseWriter, r *http.Request) {
	path := d.Config.Server.Directory
	max_results := 100 // TODO: temporary.
	finfo_list, err := readIndexFile(path, max_results)
	if err != nil {
		http.Error(w, "Could not read Index File", 500)
	}

	t, _ := template.ParseFiles("list.html")
	t.Execute(w, finfo_list)
}

func (d DevStashHTTP) showHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/s/"):]
	filename := filepath.Join(d.Config.Server.Directory, makeHashedDirName(hash))

	finfo, err := searchFile(filename, d.Config.Server.Directory)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	if strings.HasPrefix(finfo.FileType, "text/") {
		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		fmt.Fprintf(w, "<html><body><pre>%s</pre></body></html>", buf) // TODO: use template with content-type
	} else {
		fmt.Fprintf(w, "<div>Could not draw. The contents-type of %s is %s</div>", hash, finfo.FileType)
	}
}

func (d DevStashHTTP) StartHTTPServer() {
	port := d.Config.Server.Port

	http.HandleFunc("/p", d.storeHandler)
	http.HandleFunc("/", d.listHandler)
	http.HandleFunc("/s/", d.showHandler)

	addr := ":" + strconv.Itoa(port)
	fmt.Println("Start at " + addr)
	if d.Config.Server.Ssl == true {
		fmt.Println(http.ListenAndServeTLS(addr, d.Config.Server.CertFile, d.Config.Server.KeyFile, nil))
	} else {
		fmt.Println(http.ListenAndServe(addr, nil))
	}
}
