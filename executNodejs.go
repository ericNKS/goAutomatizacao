package main

import (
	"fmt"
	"os"
	"os/exec"
)

func ExecuteNodeJs(dir string) {
	var userPath, _ = os.UserHomeDir()
	cmd := exec.Command("node", userPath+"\\"+dir)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	fmt.Println("Nodejs run")
}
