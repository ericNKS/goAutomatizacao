package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	//"bytes"
	"encoding/json"
	"strconv"

	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Gabinete struct {
	name string
	dir  string
	port string
}

type Site struct {
	ID      int    `json:"id"`
	SitNome string `json:"sitNome"`
	SitDns  string `json:"sitDns"`
}

type Response struct {
	Data []Site `json:"data"`
}

var gabinetes = []Gabinete{
	//{"Teste", "plt-0025", "4025"},
}

// Function to load .env
var err = godotenv.Load(".env")

func VerifyAndRun() {
	for true {
		if !mongoIsRunning() {
			startMongo()
		}

		// obter os paineis
		url := "https://pltmax.com/s/get-painel"

		resp, err := http.Get(url + "?auth=p0l1t1M@X")
		if err != nil {
			fmt.Println("Erro ao fazer a solicitação GET:", err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Erro ao ler o corpo da resposta:", err)
		}

		var responseData Response
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			fmt.Println("Erro ao desserializar JSON:", err)

		}

		paineis := responseData.Data

		/**/
		for i := range paineis {

			dirPainel := paineis[i].SitDns
			if dirPainel == "" {
				fmt.Println("!vazio!")
				dirPainel = strconv.Itoa(paineis[i].ID)
			}

			dirPainel = "plt-" + dirPainel
			var portPainel = fmt.Sprintf("%02d", paineis[i].ID)
			portPainel = "40" + portPainel

			newSite := Gabinete{name: paineis[i].SitNome, dir: dirPainel, port: portPainel}

			gabinetes = append(gabinetes, newSite)

		}

		/**/
		for i := range gabinetes {
			if !IsGabOn(gabinetes[i].port) {
				fmt.Println("isGabOn", gabinetes[i], IsGabOn(gabinetes[i].port))
				AddDirGab(gabinetes[i])
				//initServer(gabinetes[i].dir)
				time.Sleep(50 * time.Second)
			}
		}

		time.Sleep(60 * time.Second)
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
}
func mongoIsRunning() bool {
	cmd := exec.Command("sc", "query", "mongodb")
	// Inicia o comando
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return false
	}
	return strings.Contains(string(out), "RUNNING")
}

func initServer(dir string) bool {
	if err != nil {
		return false
	}
	nodeApplication := os.Getenv("APPLICATION_DIR") + "\\" + dir
	cmd := exec.Command("cmd", "/c", "start", "cmd", "/c", "cd /d "+nodeApplication+" && node index && exit")
	// Inicia o comando
	if err := cmd.Start(); err != nil {
		return false
	} else {
		return true
	}
}
func createFile(nameFile string, fullDir string, context string) bool {
	if _, err := os.Stat(fullDir); os.IsNotExist(err) {
		if os.Mkdir(fullDir, 0755) != nil {
			return false
		}

		// Creating index.js
		f, err := os.Create(filepath.Join(fullDir, nameFile))

		if err != nil {
			fmt.Println(err)
			return false
		}
		// Continue

		_, err = f.WriteString(context)
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
		fmt.Println("The provided directory in ", fullDir, "exists")
		return false
	}
}
func RunGab(c *gin.Context) {
	dir := c.Param("dir")
	name := c.Param("name")
	port := c.Param("port")

	gab := isGab(port)
	if gab.port != port && gab.dir != dir {
		fullDir := os.Getenv("APPLICATION_DIR") + "\\" + dir
		scriptToIndexJs := `
			// define a porta padrao
			const PORT = ` + port + `;
			const PORTA = "` + port + `";
			const NOME = "` + name + `";
			const gabineteID = ` + strings.Replace(port[1:], "0", "", -1) + `;
			const banco = "` + dir + `";			
			
			/* FUNCAO IMPORTADA */
			const pltFunctions = require("../functions");
			console.log( 'PORT', PORT )
			const classPltFunctions = new pltFunctions( PORT, PORTA, NOME, gabineteID, banco );
			console.log( classPltFunctions )
		`
		if createFile("index.js", fullDir, scriptToIndexJs) {
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

func AddDirGab(gab Gabinete) {

	//fmt.Println( "AddDirGab > ", gab )

	var dir = gab.dir
	var name = gab.name
	var port = gab.port

	//gab := isGab(port)

	fullDir := os.Getenv("APPLICATION_DIR") + "\\" + dir
	fmt.Println("fullDir: ", fullDir)

	if _, err := os.Stat(fullDir); os.IsNotExist(err) {
		fmt.Println("DIR existe! ")
		scriptToIndexJs := `
			// define a porta padrao
			const PORT = ` + port + `;
			const PORTA = "` + port + `";
			const NOME = "` + name + `";
			const gabineteID = ` + strings.Replace(port[1:], "0", "", -1) + `;
			const banco = "` + dir + `";			
			
			/* FUNCAO IMPORTADA */
			const pltFunctions = require("../functions");
			console.log( 'PORT', PORT )
			const classPltFunctions = new pltFunctions( PORT, PORTA, NOME, gabineteID, banco );
			console.log( classPltFunctions )
		`
		if createFile("index.js", fullDir, scriptToIndexJs) {
			gab.port = port
			gab.dir = dir
			//gabinetes = append(gabinetes, *gab)
		} else {
			//c.JSON(http.StatusBadRequest, gin.H{"error": "gabinete or port invalid"})
			//c.Abort()
		}

	} else {
		fmt.Println("Não existe DIR!! ")
		//c.JSON(http.StatusBadRequest, gin.H{"error": "gabinete or port invalid"})
		//c.Abort()
	}
	if gab.dir == dir && gab.port == port {
		if !IsGabOn(gab.port) {
			// Inicia o servidor
			if initServer(gab.dir) {
				//c.JSON(http.StatusOK, gin.H{"success": "Server is running"})
				fmt.Println("Server is running")
			} else {
				//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting command"})
			}
		} else {
			//c.JSON(http.StatusConflict, gin.H{"error": "Server already running"})
			fmt.Println("Server already running")
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

	//fmt.Println("teste", port)
	if err != nil {
		return false
	}
	if isGab(port).port == "" {

		return false
	}
	conn.Close()
	return true
}
