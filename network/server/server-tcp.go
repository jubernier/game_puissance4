package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"network"
	"strconv"
	"sync"
)

// Server : Structure de données avec tout ce qu'il nous faut !
type Server struct {
	clients  []net.Conn // Tableau de connexions qui représentent les clients
	listener net.Listener

	writeChans []chan string // Tunnel d'écriture
	readChans  []chan string // Tunnel de lecture

	countTokens     int        // Compteur utilisé plusieurs fois dans le programme
	mutexCountRC    sync.Mutex // Mutex permettant de bloquer certaines parties du code
	finishTimes     [2]string  // Tableau des temps
	runnerPositions [2]int     // Tableau des runners sélectionnés
}

// newServer : Creating a new server with clients, listener, writeChans, readChans
func newServer() Server {
	listener, _ := net.Listen("tcp", ":8080")
	return Server{
		[]net.Conn{},
		listener,
		[]chan string{},
		[]chan string{},
		0,
		sync.Mutex{},
		[2]string{},
		[2]int{0, 1},
	}
}

// acceptClients : Accepte 2 connexions de clients
func (s *Server) acceptClients() {
	for i := 0; i < 2; i++ {
		conn, err := s.listener.Accept() // Accepte la connection du client
		if err != nil {                  // La méthode Accept est bloquante, elle ne passe à l'étape suivante (ligne suivante) que lorsqu'une connection vient d'être acceptée
			log.Fatal(err)
		}
		s.clients = append(s.clients, conn) // Ajout de la connection de ce nouveau client à la liste des clients

		// Créer les canaux de communication avec ce client
		// Enregistre ces canaux dans un tableau répertoriant tous les canaux
		s.writeChans = append(s.writeChans, make(chan string, 1))
		s.readChans = append(s.readChans, make(chan string, 1))

		// Initialiser les goroutines de communication de ce client avec la connection et le channel initialisé précédement
		go network.ReadFromNetWork(bufio.NewReader(conn), s.readChans[i])
		go network.WriteFromNetWork(bufio.NewWriter(conn), s.writeChans[i])

		s.writeChans[i] <- network.CLIENT_NUMBER + fmt.Sprint(i) // Envoyer l'ID au client qui vient de se connecter
		s.sendToAll(network.CLIENTS_IN_QUEUE + fmt.Sprint(i+1))  // Envoie à tout le monde le nombre de clients actuellement connectés

		log.Println("Client ", i, "connected")
		log.Println("in queue: ", fmt.Sprint(i+1))
	}
}

// sendToAll : Envoyer un message à tous les clients
func (s *Server) sendToAll(message string) {
	for _, chann := range s.writeChans { // Pour tous les channels d'écriture ouverts
		chann <- message // On envoie le message dans ce canal
	}
}

// comClient : Communication avec un client
func (s *Server) comClient(clientID int) {
	var message string
	for {

		// On lit le message du channel du client
		// On rappelle que cette fonction est lancée comme go routine et que le serveur lance cette fonction pour chaque client
		select {
		case message = <-s.readChans[clientID]:
			if message[:1] == network.CLIENT_CONNECTED {
				log.Println("client ", clientID, "connecté")
			}
			if message[:1] == network.CLIENT_CHOOSE_TOKEN {
				log.Println("le client", clientID, "a choisis son personage", message[1])
				s.writeChans[otherClient(clientID)] <- network.CLIENT_CHOOSE_TOKEN + strconv.Itoa(clientID) + string(message[1])
			}
			if message[:1] == network.TOKEN_POSITION {
				log.Println(message[1:])
				s.writeChans[otherClient(clientID)] <- network.TOKEN_POSITION + message[1:]
			}
			if message[:1] == network.CLIENT_TOKEN_PLAY {
				s.writeChans[otherClient(clientID)] <- network.CLIENT_TOKEN_PLAY + message[1:3]
			}
			if message[:1] == network.TOKEN_CHOICE_POSITION {
				s.writeChans[otherClient(clientID)] <- network.TOKEN_CHOICE_POSITION + message[1:]
			}
			if message[:1] == network.CLIENT_REMOVE_TOKEN {
				s.writeChans[otherClient(clientID)] <- network.CLIENT_REMOVE_TOKEN + message[1:]
			}
		}
	}
}

func main() {
	var server = newServer() // Création du serveur
	log.Println("Listening for connections")
	server.acceptClients() // Attend que les 2 clients soient connectés

	server.writeChans[0] <- network.ISHOST
	server.sendToAll(network.ALL_CONNECTED) // Envoie à tous les clients l'information que tous les clients sont connectés
	log.Println("All clients connected")
	defer server.listener.Close()         // Fermer la connection à la fin du programme
	for i, conn := range server.clients { // Pour tous les clients connectés
		defer conn.Close()
		go server.comClient(i) // On initialise une connection avec ce client
	}

	for {
	} // Pour empêcher le serveur de s'éteindre
}

func otherClient(id int) int {
	if id == 0 {
		return 1
	} else if id == 1 {
		return 0
	}
	return -1
}
