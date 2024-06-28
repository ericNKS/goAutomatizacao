package main

import (
	"fmt"
	"net"
	"os/exec"
)

func ExecuteNodeJs(nodeApplication string, nodePort string) bool {

	if !isPortInUse(nodePort) {
		cmd := exec.Command("cmd", "/c", "start", "cmd", "/k", "cd /d "+nodeApplication+" && node index")
		// Inicia o comando
		if err := cmd.Start(); err != nil {
			fmt.Printf("Erro ao iniciar o comando: %v\n", err)
			return false
		}
		return true
	}
	return false
}

func isPortInUse(port string) bool {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
