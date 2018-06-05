package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	goPath               string
	cmdFileName          string
	workSpaceRoot        string
	defaultWorkSpaceRoot string
	isProjectGoPath      bool
)

func init() {
	pwd, err := GetAbsDir()
	if err != nil {
		log.Println(err.Error())
		os.Exit(0)
	}
	defaultWorkSpaceRoot = pwd
	flag.StringVar(&goPath, "gopath", os.Getenv("GOPATH"), "设置项目 gopath")
	flag.StringVar(&workSpaceRoot, "workSpaceRoot", defaultWorkSpaceRoot, "设置项目根路径")
	flag.StringVar(&cmdFileName, "cmdFileName", "main.go", "设置启动文件名称")
	flag.BoolVar(&isProjectGoPath, "projectGoPath", false, "设置项目 GOPATH")
	flag.Parse()
}

// GetAbsDir 只支持 macos linux
func GetAbsDir() (string, error) {
	cmdResp := exec.Command("pwd")
	b, err := cmdResp.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func main() {
	err := Vsgoinit()
	if err != nil {
		log.Println(err.Error())
	}
}

// Vsgoinit 生成任务文件
func Vsgoinit() error {
	vscodePath := defaultWorkSpaceRoot + "/.vscode"
	if cmdFileName == "" {
		cmdFileName = "main.go"
	}
	if goPath == "" {
		goPath = os.Getenv("GOPATH")
	}
	if workSpaceRoot == "" {
		workSpaceRoot = defaultWorkSpaceRoot
	}

	if err := touchTasksJSON(vscodePath, cmdFileName, goPath, workSpaceRoot); err != nil {
		return err
	}
	if isProjectGoPath {
		if err := touchSettingsJSON(vscodePath, workSpaceRoot); err != nil {
			return err
		}
	}
	return nil
}

func touchTasksJSON(path, cmdFileName, gopath, workSpaceRoot string) error {
	err := mkdirFilePath(path)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path+"/tasks.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	content := tasksJSON{
		Version:          "0.1.0",
		IsShellCommand:   true,
		SuppressTaskName: true,
		ShowOutput:       "always",
		Task:             []tasks{},
	}
	run := tasks{
		TaskName:       "run",
		IsBuildCommand: true,
		Command:        "go",
		OSXPlatform: argsAndOpt{
			Args:    []string{"run", workSpaceRoot + "/" + cmdFileName},
			Options: map[string]map[string]string{"env": {"GOPATH": goPath}},
		},
	}

	build := tasks{
		TaskName: "build",
		Command:  "go",
		OSXPlatform: argsAndOpt{
			Args:    []string{"build", workSpaceRoot + "/" + cmdFileName},
			Options: map[string]map[string]string{"env": {"GOPATH": goPath}},
		},
	}
	content.Task = []tasks{run, build}
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}
	// fmt.Println(string(b))
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func touchSettingsJSON(path, workSpaceRoot string) error {
	err := mkdirFilePath(path)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path+"/settings.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	content := settingsJSON{
		Gopath: workSpaceRoot,
	}

	b, err := json.Marshal(content)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil

}

func mkdirFilePath(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 不存在
			return os.MkdirAll(path, os.ModePerm)
		}
		return err
	}
	return nil
}

type tasksJSON struct {
	Version          string  `json:"version"`
	IsShellCommand   bool    `json:"isShellCommand"`
	SuppressTaskName bool    `json:"suppressTaskName"`
	ShowOutput       string  `json:"showOutput"`
	Task             []tasks `json:"tasks"`
}

type tasks struct {
	TaskName       string     `json:"taskName"`
	IsBuildCommand bool       `json:"isBuildCommand"`
	Command        string     `json:"command"`
	OSXPlatform    argsAndOpt `json:"osx"`
}

type argsAndOpt struct {
	Args    []string                     `json:"args"`
	Options map[string]map[string]string `json:"options"`
}

type settingsJSON struct {
	Gopath string `json:"go.gopath"`
}
