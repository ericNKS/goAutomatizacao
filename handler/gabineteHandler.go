package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Gabinete struct {
	dir  string
	port string
}

var gabinetes = []Gabinete{
	{"golang", "8000"},
	{"golangc", "8001"},
}

// Function to load .env
var err = godotenv.Load(".env")

func VerifyAndRun() {
	for true {
		if !mongoIsRunning() {
			startMongo()
		}
		for i := range gabinetes {
			if !IsGabOn(gabinetes[i].port) {
				initServer(gabinetes[i].dir)
			}
			time.Sleep(1 * time.Second)
		}
		time.Sleep(10 * time.Second)
	}
}
func startMongo() {
	cmd := exec.Command("sc", "start", "mongodb")
	status := cmd.Run()
	if status != nil {
		log.Println("Error in starting mongoDB")
		return
	}
	log.Println("mongoDB is running")
	return
}
func mongoIsRunning() bool {
	cmd := exec.Command("sc", "query", "mongodb")
	// Inicia o comando
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "RUNNING")
}

func initServer(dir string) bool {
	if err != nil {
		return false
	}
	nodeApplication := os.Getenv("APPLICATION_DIR") + "\\" + dir
	cmd := exec.Command("cmd", "/c", "start", "cmd", "/c", "cd /d "+nodeApplication+" && node index")
	// Inicia o comando
	if err := cmd.Start(); err != nil {
		return false
	} else {
		return true
	}
}
func createDir(dir string, port string) bool {
	fullDir := os.Getenv("APPLICATION_DIR") + "\\" + dir
	if _, err := os.Stat(fullDir); os.IsNotExist(err) {
		fmt.Println(dir, "does not exist")
		if os.Mkdir(fullDir, 0755) != nil {
			return false
		}

		// Creating index.js
		f, err := os.Create(filepath.Join(fullDir, "index.js"))

		if err != nil {
			fmt.Println(err)
			return false
		}
		scriptToIndexJs := `
			// define a porta padrao
			const PORT = ` + port + `;
			const PORTA = "` + port + `";
			const NOME = "debora";
			const gabineteID = ` + strings.Replace(port[1:], "0", "", -1) + `;
			const banco = "` + dir + `";			
			
			/* FUNCAO IMPORTADA */
			const pltFunctions = require("../functions");
			console.log( 'PORT', PORT )
			const classPltFunctions = new pltFunctions( PORT, PORTA, NOME, gabineteID, banco );
			console.log( classPltFunctions )
		`

		_, err = f.WriteString(scriptToIndexJs)
		if err != nil {
			f.Close()
			return false

		}
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return false
		}
		// Return that the directory and file was created
		return true
	} else {
		fmt.Println("The provided directory named", dir, "exists")
		return false
	}
}
func RunGab(c *gin.Context) {
	dir := c.Param("dir")
	port := c.Param("port")

	gab := isGab(port)
	if gab.port != port && gab.dir != dir {
		if createDir(dir, port) {
			gab.port = port
			gab.dir = dir
			gabinetes = append(gabinetes, *gab)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "gabinete or port invalid"})
			c.Abort()
		}

	} else if (gab.port != port && gab.dir == dir) || (gab.port == port && gab.dir != dir) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gabinete or port invalid"})
		c.Abort()
	}
	if gab.dir == dir && gab.port == port {
		if !IsGabOn(gab.port) {
			// Inicia o servidor
			if initServer(gab.dir) {
				c.JSON(http.StatusOK, gin.H{"success": "Server is running"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting command"})
			}
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Server already running"})
		}
	}
}

func isGab(port string) *Gabinete {
	for i := range gabinetes {
		if gabinetes[i].port == port {
			return &gabinetes[i]
		}
	}
	return &Gabinete{}
}
func IsGabOn(port string) bool {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		return false
	}
	if isGab(port).port == "" {
		return false
	}
	conn.Close()
	return true
}
