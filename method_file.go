package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Method_file struct{}

func (m Method_file) SendStdin(config Config, tags []string, contents []byte) error {
	var err error

	u := config.parseUri()

	index_filepath := filepath.Join(u.Path, INDEX_FILE_NAME)
	m.ensureStoreDir(u.Path, index_filepath)

	filename := makeHashedFileName(contents)
	path := makeHashedDirName(filename)
	abs_path := filepath.Join(u.Path, path)

	finfo := NewFileInfo(tags, contents, filepath.Join(u.Path, path), "")

	err = m.store(contents, abs_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = m.addIndex(finfo, index_filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

func (m Method_file) SendFile(config Config, tags []string, path string) error {
	var err error
	u := config.parseUri()

	index_filepath := filepath.Join(u.Path, INDEX_FILE_NAME)
	m.ensureStoreDir(u.Path, index_filepath)

	fmt.Println("file: not implemented yet") // FIXME

	return err
}

func (m Method_file) List(config Config, condition string, max_results int) error {
	u := config.parseUri()

	finfo_list, err := readIndexFile(u.Path, max_results)
	for _, f := range finfo_list {
		fmt.Println(f.MakeListString("file"))
	}

	return err
}

// store contents to path
func (m Method_file) store(contents []byte, path string) error {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}

	err := ioutil.WriteFile(path, []byte(contents), 0600)
	return err
}

// add file info to index file (not FM-index)
func (m Method_file) addIndex(finfo FileInfo, index_path string) error {
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

func (m Method_file) ensureStoreDir(path string, index_filepath string) {
	// checking store dir. creates if not exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(path + " is not exists. Creating...")
		err := os.MkdirAll(path, 0700)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// create empty index file
		_, err = os.Create(index_filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}
