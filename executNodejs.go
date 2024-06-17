package main

import (
	"fmt"
	"os/exec"
)

func ExecuteNodeJs(dir string) {
	cmd := exec.Command("cmd", "/c", "start", "cmd", "/k", "node", dir)
	if cmd.Run() != nil {
		fmt.Printf("Erro ao iniciar o comando: \n")
		panic(cmd)
	}
}
