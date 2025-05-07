package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {

	// 命令行参数 Path, 缺省值为当前工作目录
	path := flag.String("path", ".", "path to the minecraft server executable directory")
	flag.Parse()

	absPath, err := filepath.Abs(*path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("path:", "Starting server at ", absPath)

	// 根据当前的OS来决定启动命令
	var cmd *exec.Cmd
	// fmt.Println(runtime.GOOS)
	if runtime.GOOS == "windows" {
		cmd = exec.Command(absPath + "\\bedrock_server.exe")
	} else if runtime.GOOS == "linux" {

	} else {
		log.Fatal("Unsupported platform", runtime.GOOS)
	}

	cmd.Dir = absPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Minecraft server started")

	err = cmd.Wait()
	if err != nil {
		log.Fatal("Minecraft server exited with error: ", err)
	}

	fmt.Println("Minecraft server exited")
}
