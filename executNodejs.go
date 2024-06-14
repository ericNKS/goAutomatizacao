package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ExecuteNodeJs(dir string) {
	userPath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// Construir o caminho completo para o script Node.js
	scriptPath := filepath.Join(userPath, dir)

	// Comando para abrir o prompt de comando e executar o script Node.js
	cmd := exec.Command("cmd", "/C", "start", "cmd", "/K", "node", scriptPath)

	// Definir o diretório de trabalho
	cmd.Dir = filepath.Join(userPath, dir)

	// Iniciar o comando
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Erro ao iniciar o comando: %s\n", err)
		return
	}

	fmt.Println("Node.js script iniciado em uma nova janela do prompt")
}
