package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

const (
	templateFilename = "./.commit_template"
	GitDir           = ".git"
	HooksDir         = "hooks"
	HookFile         = "prepare-commit-msg"
)

// standard input scanner
var sc = bufio.NewScanner(os.Stdin)

// main function
func main() {
	fmt.Println("create template")
	createTemplate()
	fmt.Println("create prepare")
	createPrepare()
}

// search git dir
func findGitdir(path string) (string, error) {
	searchPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	} else if searchPath == "/" {
		return "", errors.New("owari")
	}

	files, err := ioutil.ReadDir(searchPath)
	if err != nil {
		return "", err
	}

	for _, val := range files {
		if val.IsDir() && val.Name() == GitDir {
			fullPath := filepath.Join(searchPath, val.Name())
			return fullPath, nil
		}
	}
	return findGitdir(filepath.Join(searchPath, "../"))
}

// make prepare file
func createPrepare() {
	path, err := findGitdir("./")
	if err != nil {
		fmt.Println(err)
		fmt.Println("End App")
		return
	}
	data := `#!/bin/sh
## prepare-commit-msg
branch=$(git branch | grep "*" | awk '{print $2}' | sed -e 's/[\/_]/ /g')
perl -i.bak -ne "s/{branch}/$branch/g; print" "$1"`
	fPath := filepath.Join(path, HooksDir, HookFile)
	createHook(fPath, data)
	if err = os.Chmod(fPath, 0777); err != nil {
		log.Println(err)
	}
}

// create prepare file
func createHook(fPath string, data string) {
	file, err := os.Create(fPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	if _, err = file.WriteString(data); err != nil {
		log.Fatal(err)
	}
}

// create home dir
func createHomedir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func createTemplate() {

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

	vals := []string{"config", "commit.template", cmd}
	if err := exec.Command("git", vals...).Run(); err != nil {
		fmt.Printf("Val:%v\nError:%v\n", vals, err)
	}
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
	defer func() {
		defer f.Close()
		if err = os.Chmod(fPath, 0777); err != nil {
			log.Println(err)
		}
	}()
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
