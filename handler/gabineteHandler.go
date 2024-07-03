package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net"
	"net/http"
	"os"
	"os/exec"
)

type Gabinete struct {
	name string
	dir  string
	port string
}

var gabinetes = []Gabinete{
	{"Teste 1", "golang", "8000"},
	{"Teste 2", "golangc", "8001"},
}

func NewGab(name string, dir string, port string) *Gabinete {
	gab := Gabinete{name: name, dir: dir, port: port}
	gabinetes = append(gabinetes, gab)
	return &gab
}

func RunGab(c *gin.Context) {
	dir := c.Param("dir")
	port := c.Param("port")

	// Function to load .env
	err := godotenv.Load(".env")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	gab := isGab(port)
	if gab.port != "" && gab.dir == dir {
		nodeApplication := os.Getenv("APPLICATION_DIR") + "\\" + gab.dir
		if !IsGabOn(port) {
			cmd := exec.Command("cmd", "/c", "start", "cmd", "/k", "cd /d "+nodeApplication+" && node index")
			// Inicia o comando
			if err := cmd.Start(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting command"})
			} else {
				c.JSON(http.StatusOK, true)
			}
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Server already running"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gabinete or port invalid"})
	}
}

func isGab(port string) Gabinete {
	for i := range gabinetes {
		if gabinetes[i].port == port {
			return gabinetes[i]
		}
	}
	return Gabinete{}
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
