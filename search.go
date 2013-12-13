package main

import (
	"bitbucket.org/oov/go-shellinford/shellinford"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func EnsureFMIndex(path string) (*shellinford.FMIndex, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		of, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer of.Close()
		return shellinford.NewFMIndex(), nil
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("ERROR: Open FMIndex File")
		return nil, err
	}
	defer f.Close()
	fm, err := shellinford.OpenFMIndex(f)
	if err != nil {
		fmt.Printf("ERROR: OpenFMIndex")
		return nil, err
	}

	return fm, nil
}

func MakeIndex(config Config) error {
	u, err := url.Parse(config.Default.Uri)
	if err != nil {
		return err
	}
	index_filepath := filepath.Join(u.Path, INDEX_FILE_NAME)
	fm_index_filepath := filepath.Join(u.Path, FM_INDEX_FILE_NAME)
	fm, err := EnsureFMIndex(fm_index_filepath)
	if err != nil {
		return err
	}

	// read all contents at once
	contents, err := ioutil.ReadFile(index_filepath)
	if err != nil {
		fmt.Printf("ERROR: Read index file")
		return err
	}
	// split by new line
	lines := strings.Split(string(contents), "\n")

	fmt.Printf("Index creating ")

	for _, l := range lines {
		finfo := ReadFileInfo(l)
		add(fm, config, finfo.Path)
	}

	of, err := os.OpenFile(fm_index_filepath, os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("ERROR: Open FMIndex File at MakeIndex")
		return err
	}

	fmt.Println(" done")

	defer of.Close()
	err = fm.Write(of)
	if err != nil {
		fmt.Printf("ERROR: Write FMIndex File at MakeIndex")
		return err
	}

	return nil
}

func add(fm *shellinford.FMIndex, config Config, filepath string) error {
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	fm.Add(contents)

	fm.Build(0, 1)
	fmt.Printf(".")

	return nil
}

func Search(config Config, searchword string) error {
	u, err := url.Parse(config.Default.Uri)
	if err != nil {
		return err
	}
	fm_index_filepath := filepath.Join(u.Path, FM_INDEX_FILE_NAME)
	fm, err := EnsureFMIndex(fm_index_filepath)
	if err != nil {
		return err
	}

	st := time.Now()
	found := fm.Search([]byte(searchword))
	fmt.Printf("Keyword \"%s\" found %d entries.\n", searchword, len(found), time.Now().Sub(st).Seconds())
	for k, v := range found {
		fmt.Println(string(fm.Document(k)), v, "Hit(s)")
	}

	return nil
}
