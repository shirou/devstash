package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Method_ssh struct{}

func (m Method_ssh) SendStdin(config Config, tags []string, contents []byte) error {
	u := config.parseUri()

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

	filename := makeHashedFileName(contents)
	path := makeHashedDirName(filename)
	remote_file := filepath.Join(u.Path, path)

	finfo := NewFileInfo(tags, contents, filepath.Join(u.Path, path), "")

	server_name := u.Host
	// TODO
	//	server_name := u.User.Username() + u.Host

	m.sshSend(server_name, f.Name(), remote_file, u.Path, finfo)

	return nil
}

// Internal implementation of two ssh send functions.
// send contents using /bin/cat and add Index to index file on remote
// these two job run on one ssh connection to reduce time.
func (m Method_ssh) sshSend(hostname string, local string, remote string, remote_root string, finfo FileInfo) error {

	dir_path := filepath.Dir(remote)
	index_filepath := filepath.Join(remote_root, INDEX_FILE_NAME)

	// mkdir, cat contents and append indexfile with one ssh connection
	remote_cat := fmt.Sprintf("mkdir -p %s && /bin/cat > %s && echo %s >> %s",
		dir_path, remote, strings.Replace(finfo.MakeIndexFormat(), "\t", `\\t`, -1), index_filepath)

	// exec ssh
	cat := exec.Command("cat", local)
	ssh := exec.Command("ssh", "-C", hostname, remote_cat)

	pipe, err := ssh.StdinPipe()
	if err != nil {
		return err
	}
	cat.Stdout = pipe

	cat.Start()
	ssh.Start()
	cat.Wait()

	return nil
}

func (m Method_ssh) SendFile(config Config, tags []string, path string) error {
	u := config.parseUri()

	fmt.Println("ssh: not implemented yet" + u.Path) // FIXME
	return nil
}

// get indexfile and show lists
// TODO: but there is no head lines is this necessary?
func (m Method_ssh) List(config Config, condition string, max_results int) error {
	u := config.parseUri()

	server_name := u.Host
	// TODO:
	//	if u.User.Username() != "" {
	//		server_name = u.User.Username() + u.Host
	//	}

	index_filepath := filepath.Join(u.Path, INDEX_FILE_NAME)

	cmd := fmt.Sprintf("tail -%d %s | sort -r ", max_results, index_filepath)
	ssh := exec.Command("ssh", server_name, cmd)
	out, err := ssh.Output()

	if err != nil {
		fmt.Println("ERROR: ssh List")
		return err
	}

	for _, l := range strings.Split(string(out), "\n") {
		if len(l) < 3 {
			break
		}
		f := ReadFileInfo(l)
		fmt.Println(f.MakeListString("ssh"))
	}
	return nil
}
