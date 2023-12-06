package main

import (
	"log"
	"net"
	"sync"
)

func handleServer(c net.Conn) {
	var logLock sync.Mutex
	inBuf := make([]byte, 22)

	for {
		logLock.Lock()
		_, err := c.Read(inBuf)
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		logLock.Unlock()
	log.Println("Received from server:", string(inBuf))

	
		
	}
}

// Création, paramétrage et lancement du jeu.
func main() {
	//connection
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Dial error:", err)
		return
	}
	defer conn.Close()

	log.Println("Connected to the server")

	handleServer(conn)


}