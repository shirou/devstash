package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// file information. This will be added to index file as LTSV
type FileInfo struct {
	Created        time.Time
	Tags           []string
	Path           string // stored abs filepath
	OrigPath       string // if file
	SenderAddr     string // sender host address if http or ssh
	SenderHostname string // sender host name if http or ssh
	Head           string // first line of contents.
	FileType       string // MIME filetype
}

// Make one line string from a FileInfo to append a index file
// Line is Labeled Tab-separated Values (LTSV)
// see http://ltsv.org
func (finfo FileInfo) MakeIndexFormat() string {
	buf := []string{
		"time:" + finfo.Created.Format(time.RFC3339),
		"tags:" + strings.Join(finfo.Tags, ","),
		"path:" + finfo.Path,
		"orig:" + finfo.OrigPath,
		"addr:" + finfo.SenderAddr,
		"host:" + finfo.SenderHostname,
		"head:" + finfo.Head,
		"filetype:" + finfo.FileType,
	}

	return strings.Join(buf, "\t") // return as LTSV
}

// make string which is used by list
func (finfo FileInfo) MakeListString(method string) string {
	//	MAX_HEAD := 10
	format := "%s\t%s\t%s\t%s"

	time := finfo.Created.Format(time.RFC3339)[:19] // only needs sec
	tag := strings.Join(finfo.Tags, ",")
	path := filepath.Base(finfo.Path)[:FILENAME_LEN]
	head := finfo.Head
	return fmt.Sprintf(format, time, tag, path, head)
}

func NewFileInfo(tags []string, contents []byte, stored string, orig string) FileInfo {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Could not get hostname: %v\n", err)
		os.Exit(1)
	}

	addrs, err := net.LookupHost(hostname)
	if err != nil {
		fmt.Printf("Could not get ipaddress: %v\n", err)
		os.Exit(1)
	}

	return newFileInfoImpl(tags, contents, stored, orig, hostname, addrs[1])
}


func NewFileInfoWithAddr(tags []string, contents []byte, stored string, orig string, hostname string, ipaddress string) FileInfo {
	return newFileInfoImpl(tags, contents, stored, orig, hostname, ipaddress)
}


func newFileInfoImpl(tags []string, contents []byte, stored string, orig string, hostname string, ipaddress string) FileInfo {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Could not get hostname: %v\n", err)
		os.Exit(1)
	}

	addrs, err := net.LookupHost(hostname)
	if err != nil {
		fmt.Printf("Could not get ipaddress: %v\n", err)
		os.Exit(1)
	}

	filetype := http.DetectContentType(contents)

	head := ""
	if strings.HasPrefix(filetype, "text/") {
		head = strings.Split(string(contents), "\n")[0] // get first line if text
	}

	return FileInfo{
		time.Now(),
		tags,
		stored,
		orig,
		addrs[0], // only first address
		hostname,
		head,
		filetype,
	}
}


func ReadFileInfo(line string) FileInfo {
	f := strings.Split(strings.TrimSpace(line), "\t")

	finfo := FileInfo{}
	for _, l := range f {
		lv := strings.Split(l, ":")
		label := lv[0]
		value := strings.Join(lv[1:], ":")

		switch label {
		case "time":
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				fmt.Println(err)
				return finfo
			}
			finfo.Created = t
		case "tags":
			finfo.Tags = strings.Split(value, ",")
		case "path":
			finfo.Path = value
		case "orig":
			finfo.OrigPath = value
		case "addr":
			finfo.SenderAddr = value
		case "host":
			finfo.SenderHostname = value
		case "head":
			finfo.Head = value
		case "filetype":
			finfo.FileType = value
		}
	}

	return finfo
}

func (finfo FileInfo) Basename() string {
	return filepath.Base(finfo.Path)
}
func (finfo FileInfo) LinkPath() string {
	return "/s/" + filepath.Base(finfo.Path)
}


// add file info to index file (not FM-index)
func (finfo FileInfo) addIndex(index_path string) error {
	f, err := os.OpenFile(index_path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	append_line := finfo.MakeIndexFormat() + "\n"

	if _, err = f.WriteString(append_line); err != nil {
		return err
	}
	return nil
}

