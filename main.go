package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"strings"
)

const version = "0.0.1"
const INDEX_FILE_NAME = "stash.idx"
const FM_INDEX_FILE_NAME = "fmstatsh.idx"

type Method interface {
	SendStdin(Config, []string, []byte) error
	SendFile(Config, []string, string) error
	List(Config, string, int) error
}

// config file
type Config struct {
	Default struct {
		Uri string
	}
	Server struct {
		Port      int
		Directory string
		Ssl       bool
		CertFile  string
		KeyFile   string
	}
}

func (c Config) parseUri() url.URL {
	// parse
	u, err := url.Parse(c.Default.Uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return *u
}

func LoadConfig(config_path string) (Config, error) {
	var err error
	var cfg Config

	// replace ~ to home directory
	usr, _ := user.Current()
	dir := usr.HomeDir
	config_path = strings.Replace(config_path, "~", dir, 1)

	err = gcfg.ReadFileInto(&cfg, config_path)
	return cfg, err
}

func main() {
	var filename string
	var server_mode bool
	var list_mode bool
	var search_mode string
	var index_create bool
	var config_path string
	var max_results int
	var tags_arg string

	flag.StringVar(&filename, "f", "", "Filename to store")
	flag.StringVar(&search_mode, "s", "", "Search string")
	flag.StringVar(&config_path, "c", "~/.devstash.cfg", "Config path")
	flag.BoolVar(&server_mode, "server", false, "run as http server")
	flag.BoolVar(&list_mode, "l", false, "List stash")
	flag.BoolVar(&index_create, "make-index", false, "Create Index")
	flag.IntVar(&max_results, "n", 100, "Max results in list")
	flag.StringVar(&tags_arg, "t", "", "Tags (ex: a,b,c)")

	flag.Parse()

	// Load config file
	conf, err := LoadConfig(config_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// make tags
	tags := strings.Split(tags_arg, ",")

	// create method which is specified default section in config file
	var method Method
	if strings.HasPrefix(conf.Default.Uri, "ssh://") {
		method = Method_ssh{}
	} else if strings.HasPrefix(conf.Default.Uri, "http://") || strings.HasPrefix(conf.Default.Uri, "https://") {
		fmt.Println("http")
	} else if strings.HasPrefix(conf.Default.Uri, "file://") {
		method = Method_file{}
	} else {
		fmt.Println("Unknown method:" + conf.Default.Uri)
		os.Exit(1)
	}

	if server_mode == true {
		d := DevStashHTTP{conf}
		d.StartHTTPServer()
		os.Exit(0)
	}
	if list_mode == true {
		method.List(conf, "", max_results)
		os.Exit(0)
	}
	if index_create == true {
		err := MakeIndex(conf)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}
	if search_mode != "" {
		Search(conf, search_mode)
		os.Exit(0)
	}

	if filename == "" {
		contents, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
		}
		method.SendStdin(conf, tags, contents)
	} else {
		method.SendFile(conf, tags, filename)
	}

}
