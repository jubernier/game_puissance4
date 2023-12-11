package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"

	"puissancequatre/network"

	"github.com/hajimehoshi/ebiten/v2"
)

// Structure de données pour représenter l'état courant du jeu.
type game struct {
	gameState     int
	stateFrame    int
	grid          [globalNumTilesX][globalNumTilesY]int
	p1Color       int
	p2Color       int
	turn          int
	tokenPosition int
	result        int
	readChan      chan string
	writeChan     chan string
	clientId      int
}

// Constantes pour représenter la séquence de jeu actuelle (écran titre,
// écran de sélection des couleurs, jeu, écran de résultats).
const (
	titleState int = iota
	colorSelectState
	playState
	resultState
)

// Constantes pour représenter les pions dans la grille de puissance 4
// (absence de pion, pion du joueur 1, pion du joueur 2).
const (
	noToken int = iota
	p1Token
	p2Token
)

// Constantes pour représenter le tour de jeu (joueur 1 ou joueur 2).
const (
	p1Turn int = iota
	p2Turn
)

// Constantes pour représenter le résultat d'une partie (égalité si
// la grille a été remplie sans qu'un joueur n'ait gagné, joueur 1
// gagnant ou joueur 2 gagnant).
const (
	equality int = iota
	p1wins
	p2wins
)

// Remise à 0 du jeu pour recommencer une partie. Le joueur qui a
// perdu la dernière partie commence.
func (g *game) reset() {
	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {
			g.grid[x][y] = noToken
		}
	}
}

func InitGame(ip, port string) (g game) {
	// Open connection
	log.Println(ip + ":" + port)
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	// Initialisation du channel de communication
	g.readChan = make(chan string, 1)
	g.writeChan = make(chan string, 1)

	// Goroutine écoutant permettant le lire en double sur un reader initialisé avec la connection
	go network.server.ReadFromNetWork(bufio.NewReader(conn), g.readChan)
	go network.server.WriteFromNetWork(bufio.NewWriter(conn), g.writeChan)

	var message = <-g.readChan
	if message[:1] == network.CLIENT_NUMBER {
		var idFromServ, _ = strconv.Atoi(message[1:])
		g.clientId = idFromServ
		ebiten.SetWindowTitle("BUT2 année 2022-2023, R3.05 Programmation système, clientID: " + fmt.Sprint(g.clientId))
	}

	g.writeChan <- network.CLIENT_CONNECTED

	return g
}
