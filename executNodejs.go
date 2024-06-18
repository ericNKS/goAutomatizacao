package main

import (
	"fmt"
	"os/exec"
)

func ExecuteNodeJs(dir string) {
	//cmd := exec.Command("cmd", "/c", "start", "cmd", "/k", "node", dir)
	cmd := exec.Command("cmd", "/c", "start", "cmd", "/k", "cd /d "+dir+" && node index")
	if cmd.Run() != nil {
		fmt.Printf("Erro ao iniciar o comando: \n")
		panic(cmd)
	}
}
