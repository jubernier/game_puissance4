package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"network"
	"strconv"
)

// Mise à jour de l'état du jeu en fonction des entrées au clavier.
func (g *Game) Update() error {

	g.stateFrame++

	switch g.gameState {
	case titleState:
		if g.titleUpdate() {
			g.gameState++
		}
	case colorSelectState:
		if g.colorSelectUpdate() {
			if !(g.clientId == 0) {
				g.turn = p2Turn
			}
			g.gameState++
		}
	case playState:
		var lastXPositionPlayed int
		var lastYPositionPlayed int
		if g.turn == p1Turn {
			g.tokenPosUpdate()
			lastXPositionPlayed, lastYPositionPlayed = g.p1Update()
		} else {
			lastXPositionPlayed, lastYPositionPlayed = g.p2Update()
		}
		if lastXPositionPlayed >= 0 {
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished {
				g.result = result
				g.gameState++
			}
		}
	case resultState:
		if g.resultUpdate() {
			g.reset()
			g.gameState = playState
		}
	}

	return nil
}

// Mise à jour de l'état du jeu à l'écran titre.
func (g *Game) titleUpdate() bool {
	g.stateFrame = g.stateFrame % globalBlinkDuration
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *Game) colorSelectUpdate() bool {
	col := g.p1Color % globalNumColorCol
	line := g.p1Color / globalNumColorLine

	if !g.p1ChooseToken {
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			col = (col + 1) % globalNumColorCol
			//change = true
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			col = (col - 1 + globalNumColorCol) % globalNumColorCol
			//change = true
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			line = (line + 1) % globalNumColorLine
			//change = true
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			line = (line - 1 + globalNumColorLine) % globalNumColorLine
			//change = true
		}
		g.p1Color = line*globalNumColorLine + col

		if g.p1Change != g.p1Color {
			moveMessage := network.TOKEN_CHOICE_POSITION + strconv.Itoa(g.p1Color)
			g.p1Change = g.p1Color
			g.writeChan <- moveMessage
		}
		select {
		case message := <-g.readChan:
			if message[:1] == network.CLIENT_REMOVE_TOKEN {
				g.p2ChooseToken = false
				log.Println("le p2 à désectionner son token", g.clientId, message[1:])
			}
			if message[:1] == network.CLIENT_CHOOSE_TOKEN {
				g.p2ChooseToken = true
				g.p2Color, _ = strconv.Atoi(string(message[2]))
				if g.p1ChooseToken {
					return true
				}
			}
			if message[:1] == network.TOKEN_CHOICE_POSITION {
				pos, _ := strconv.Atoi(message[1:])
				//log.Println("POSITION RECUE: ", message[1:])
				g.p2Color = pos
				if g.p1Color == g.p2Color {
					g.p1Color = (g.p1Color + 1) % globalNumColor
					moveMessage := network.TOKEN_CHOICE_POSITION + strconv.Itoa(g.p1Color)
					//log.Println("SEND TO SERVER: " + moveMessage)
					g.writeChan <- moveMessage
				}
			}

		default:
		}
	}
	/*
		println(g.p1Color)
		if change {
			var msg string = network.TOKEN_CHOICE_POSITION + strconv.Itoa(g.p1Color)

			g.writeChan <- msg
		}
	*/

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.writeChan <- network.CLIENT_CHOOSE_TOKEN + strconv.Itoa(g.p1Color)
		g.p1ChooseToken = true
		if g.p2ChooseToken == true {
			log.Println("All token choosen")
			return true
		}
		return false
		//g.writeChan <- network.TOKEN_CHOICE_POSITION + strconv.Itoa(g.clientId) + strconv.Itoa(g.p2Color)

	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.writeChan <- network.CLIENT_REMOVE_TOKEN + strconv.Itoa(g.p1Color)
		g.p1ChooseToken = false
	}
	return false

}

// Gestion de la position du prochain pion à jouer par le joueur 1.
func (g *Game) tokenPosUpdate() {

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
		g.writeChan <- network.TOKEN_POSITION + strconv.Itoa(g.tokenPosition)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
		g.writeChan <- network.TOKEN_POSITION + strconv.Itoa(g.tokenPosition)
	}

}

// Gestion du moment où le prochain pion est joué par le joueur 1.
func (g *Game) p1Update() (int, int) {

	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
			g.turn = p2Turn
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos
		}
		g.writeChan <- network.CLIENT_TOKEN_PLAY + strconv.Itoa(lastXPositionPlayed) + strconv.Itoa(lastYPositionPlayed)
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Gestion de la position du prochain pion joué par le joueur 2 et
// du moment où ce pion est joué.
func (g *Game) p2Update() (int, int) {
	var lastYpos = -1
	var lastXpos = -1
	select {
	case message := <-g.readChan:
		//log.Println("on est la", message)
		if message[:1] == network.TOKEN_POSITION {
			position, _ := strconv.Atoi(string(message[1]))
			//log.Println(position)

			g.tokenPosition = position

			return lastXpos, lastYpos
		}
		if message[:1] == network.CLIENT_TOKEN_PLAY {
			position, _ := strconv.Atoi(string(message[1]))
			updated, yPos := g.updateGrid(p2Token, position)
			for ; !updated; updated, yPos = g.updateGrid(p2Token, position) {
				position = (position + 1) % globalNumTilesX
			}
			g.turn = p1Turn
			return position, yPos
		}
	default:
	}
	return lastXpos, lastYpos
}

/*
func (g *Game) p2Update() (int, int) {
	position := rand.Intn(globalNumTilesX)
	updated, yPos := g.updateGrid(p2Token, position)
	for ; !updated; updated, yPos = g.updateGrid(p2Token, position) {
		position = (position + 1) % globalNumTilesX
	}
	g.turn = p1Turn
	return position, yPos
}*/

// Mise à jour de l'état du jeu à l'écran des résultats.
func (g Game) resultUpdate() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de la grille de jeu lorsqu'un pion est inséré dans la
// colonne de coordonnée (x) position.
func (g *Game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

// Vérification de la fin du jeu : est-ce que le dernier joueur qui
// a placé un pion gagne ? est-ce que la grille est remplie sans gagnant
// (égalité) ? ou est-ce que le jeu doit continuer ?
func (g Game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// egalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
