package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

const (
	templateFilename = "./.commit_template"
)

var sc = bufio.NewScanner(os.Stdin)

func createHomedir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func main() {
	dirPath := filepath.Join(createHomedir(), ".config/chok")
	createDirIfNeed(dirPath)
	if err := createTempFile(dirPath); err != nil {
		fmt.Println(err)
		return
	}
	cmd := filepath.Join(dirPath, templateFilename)
	systemExecIfNeed(cmd)
}
func systemExecIfNeed(cmd string) {

	vals := []string{"config", "commit_template", cmd}
	err := exec.Command("git", vals...).Run()
	fmt.Printf("Val:%v\nError:%v\n", vals, err)
}

func createDirIfNeed(path string) {
	os.MkdirAll(path, os.ModePerm)
}

func createTempFile(path string) error {
	fPath := filepath.Join(path, templateFilename)
	if Exists(fPath) {
		if err := ask("templatefile over write?", "file exists"); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed filecreate:" + fPath)
	}
	defer f.Close()
	f.WriteString("{branch}")
	return nil
}

func ask(question string, errMsg string) error {
	fmt.Println(question)
	sc.Scan()
	if sc.Text() == "y" {
		return nil
	}
	return fmt.Errorf(errMsg)
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
