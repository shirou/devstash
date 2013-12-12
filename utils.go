package main

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"time"
	"strings"
	"sort"
	"io/ioutil"
)

const FILENAME_LEN = 40

// Get sha256-hexed string from contents.
// This is used as a filename.
func makeHashedFileName(contents []byte) string {
	hasher := sha256.New()

	now := time.Now()
	now_b, _ := now.GobEncode()  // append current time to make unique
	t := hasher.Sum(append(contents, now_b...))

	return hex.EncodeToString(t)[:FILENAME_LEN]
}

// Get stored directory name from hashed-filename.
// Take 0-1 and 2-3 bytes from first and used as directory.
// ex)
//   aabbccdd  ->  /aa/bb/aabbccdd
func makeHashedDirName(filename string) string {

	s1 := filename[0:2]
	s2 := filename[2:4]

	return filepath.Join(s1, s2, filename)
}

func readIndexFile(path string, max_results int) ([]FileInfo, error) {
	// read all contents at once
	index_filepath := filepath.Join(path, INDEX_FILE_NAME)
	contents, err := ioutil.ReadFile(index_filepath)
	if err != nil {
		return []FileInfo{}, err
	}
	// split by new line
	lines := strings.Split(string(contents), "\n")
	sort.Sort(sort.Reverse(sort.StringSlice(lines)))  // reverse from latest

	max := max_results
	if len(lines) < max_results {
		max = len(lines) -1
	}

	finfo_list := []FileInfo{}
	for _, l := range lines[0:max]{
		f := ReadFileInfo(l)
		finfo_list = append(finfo_list, f)
	}

	return finfo_list, nil
}
