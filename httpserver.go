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
	"os"
	"io"
	"mime/multipart"
)

type DevStashHTTP struct {
	Config Config
}

func (d DevStashHTTP) storeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST"{
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//parse the multipart
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := r.MultipartForm

	files := m.File["devstash"]

	index_filepath := filepath.Join(d.Config.Server.Directory, INDEX_FILE_NAME)

	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dst, contents, err := d.storeMultiPartFile(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		finfo := NewFileInfoWithAddr([]string{}, contents, dst, files[i].Filename,"", r.RemoteAddr)

		err = finfo.addIndex(index_filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	fmt.Fprintf(w, "Uploaded, %q", html.EscapeString(r.URL.Path))
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

	http.HandleFunc("/", d.listHandler)
	http.HandleFunc("/p", d.storeHandler)
	http.HandleFunc("/s/", d.showHandler)

	addr := ":" + strconv.Itoa(port)
	fmt.Println("Start at " + addr)
	if d.Config.Server.Ssl == true {
		fmt.Println(http.ListenAndServeTLS(addr, d.Config.Server.CertFile, d.Config.Server.KeyFile, nil))
	} else {
		fmt.Println(http.ListenAndServe(addr, nil))
	}
}


func (d DevStashHTTP) storeMultiPartFile(file multipart.File) (string, []byte, error){
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return "", []byte{}, err
	}
	path := makeHashedDirName(makeHashedFileName(contents))
	dir := filepath.Dir(path)

	// ensure directory is created
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return "", []byte{}, err
		}
	}

	dst, err := os.Create(path)
	defer dst.Close()
	if err != nil {
		return "", []byte{}, err
	}
	//copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		return "", []byte{}, err
	}

	return path, contents, nil
}
