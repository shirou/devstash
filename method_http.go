package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Method_http struct{}

const PARAM_NAME = "File"

func (m Method_http) SendStdin(config Config, tags []string, contents []byte) error {
	var err error

	// create temp file to store contents from stdin
	f, err := ioutil.TempFile("", "devstash")
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name()) // delete temp file

	_, err = f.Write(contents)
	if err != nil {
		return err
	}

	err = m.httpPost(config, tags, f.Name())

	return err
}

func (m Method_http) SendFile(config Config, tags []string, path string) error {
	err := m.httpPost(config, tags, path)

	return err
}

func (m Method_http) List(config Config, condition string, max_results int) error {
	//	u := config.parseUri()

	fmt.Println("file: not implemented yet") // FIXME

	return nil
}

// send file implementation
func (m Method_http) httpPost(config Config, tags []string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(PARAM_NAME, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return err
	}

	u := config.parseUri()
	http.NewRequest("POST", u.String()+"/p", body)

	return nil
}
