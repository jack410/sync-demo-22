package main

import (
	"example.com/m/server"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	//启动gin服务
	go server.Run()
	cmd := startBrower()
	chSignal := listenToInterrupt()
	//如果没有值则一直等待（阻塞），直到有信号输入
	//select可以监听多个管道，只要有一个管道有信号则进行下一步
	//如果没有信号，select就等待（阻塞）
	select {
	case <-chSignal:
		//一旦有信号则关闭浏览器
		cmd.Process.Kill()
	}
}

func startBrower() *exec.Cmd {
	chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	var tmpDir string
	name, _ := ioutil.TempDir("", "lorca")
	tmpDir = name
	//删除缓存文件
	defer os.RemoveAll(tmpDir)
	fmt.Println(tmpDir)
	cmd := exec.Command(chromePath, "--app=http://localhost:27149/static/index.html", fmt.Sprintf("--user-data-dir=%s", tmpDir),
		"--no-first-run")
	cmd.Start()
	return cmd
}

//监听中断信号
func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	//signal.Notify订阅os.Interrupt信号,一旦有信号则往chSingal管道里写一个信号
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
