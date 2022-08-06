package main

import (
	"fmt"
	"github.com/andresbott/x1-acpi-manager/manager"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("this application only runs on linux")
		os.Exit(1)
	}

	if !isRoot() {
		fmt.Println("please run as root")
		os.Exit(1)
	}

	m := manager.Manager{}

	done := make(chan bool, 1)
	// handle clean shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		m.Stop()
		done <- true
	}()

	m.Start()
	<-done
	fmt.Println("clean shutdown finished")
}

func isRoot() bool {
	return os.Geteuid() == 0
}
