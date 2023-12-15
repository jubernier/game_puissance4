package main

/*
import (
	"bufio"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"net"
	"network"
	"strconv"
)

func handleServeur(ip, port string) string {
	// Open connection

	log.Println(ip + ":" + port)
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	// Initialisation du channel de communication
	var readChan = make(chan string, 1)
	var writeChan = make(chan string, 1)

	// Goroutine écoutant permettant le lire en double sur un reader initialisé avec la connection
	go network.ReadFromNetWork(bufio.NewReader(conn), readChan)
	go network.WriteFromNetWork(bufio.NewWriter(conn), writeChan)

	var message = <-readChan
	if message[:1] == network.CLIENT_NUMBER {
		var idFromServ, _ = strconv.Atoi(message[1:])
		var clientId = idFromServ
		ebiten.SetWindowTitle("clientID: " + fmt.Sprint(clientId))
	}

	writeChan <- network.CLIENT_CONNECTED
	return <-writeChan
}

func main() {
	handleServeur("127.0.0.1", "8080")
}
*/
