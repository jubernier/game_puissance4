package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strconv"
)

// Affichage des graphismes à l'écran selon l'état actuel du jeu.
func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(globalBackgroundColor)

	switch g.gameState {
	case titleState:
		g.titleDraw(screen)
	case colorSelectState:
		g.colorSelectDraw(screen)
	case playState:
		g.playDraw(screen)
	case resultState:
		g.resultDraw(screen)
	}

}

// Affichage des graphismes de l'écran titre.
func (g Game) titleDraw(screen *ebiten.Image) {
	text.Draw(screen, "Puissance 4 en réseau", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "Projet de programmation système", smallFont, 105, 190, globalTextColor)
	text.Draw(screen, "Année 2023-2024", smallFont, 210, 230, globalTextColor)
	var message = "Le nombre de joueurs connectés : " + strconv.Itoa(g.clientInQueue)
	// Affichage du nombre de joueurs connectés
	text.Draw(screen, message, smallFont, 105, 300, globalTextColor)
	if g.isReadyNextStep {
		if g.stateFrame >= globalBlinkDuration/3 {
			text.Draw(screen, "Appuyez sur entrée", smallFont, 210, 500, globalTextColor)
		}
	}
}

// Affichage des graphismes de l'écran de sélection des couleurs des joueurs.
func (g Game) colorSelectDraw(screen *ebiten.Image) {
	text.Draw(screen, "Quelle couleur pour vos pions ?", smallFont, 110, 80, globalTextColor)

	line := 0
	col := 0
	for numColor := 0; numColor < globalNumColor; numColor++ {

		xPos := (globalNumTilesX-globalNumColorCol)/2 + col
		yPos := (globalNumTilesY-globalNumColorLine)/2 + line
		// Ajout de l'affichage de la sélection du joueur 2
		if numColor == g.p2Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, p2GlobalSelectColor, true)
		}
		if numColor == g.p1Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, globalSelectColor, true)
		}
		// Changement de la couleur de sélection lorsque l'un des joueurs a choisis sa couleur
		if numColor == g.p1Color && g.p1ChooseToken {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, globalColorChoose, true)
		}
		if numColor == g.p2Color && g.p2ChooseToken {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, globalColorChoose, true)
		}
		vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2-globalCircleMargin, globalTokenColors[numColor], true)

		col++
		if col >= globalNumColorCol {
			col = 0
			line++
		}
	}
}

// Affichage des graphismes durant le jeu.
func (g Game) playDraw(screen *ebiten.Image) {
	g.drawGrid(screen)
	if g.turn == p1Turn {
		vector.DrawFilledCircle(screen, float32(globalTileSize/2+g.tokenPosition*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p1Color], true)
		// Ajout de l'information indiquant le tour de p1
		text.Draw(screen, "ton tour", smallFont, 15, 30, globalTextColor)
		// Ajout de l'affichage de p2 avec sa couleur
	} else {
		vector.DrawFilledCircle(screen, float32(globalTileSize/2+g.tokenPosition*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p2Color], true)
	}
}

// Affichage des graphismes à l'écran des résultats.
func (g Game) resultDraw(screen *ebiten.Image) {
	g.drawGrid(offScreenImage)

	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.2)
	screen.DrawImage(offScreenImage, options)

	message := "Égalité"
	if g.result == p1wins {
		message = "Gagné !"
	} else if g.result == p2wins {
		message = "Perdu…"
	}
	ebiten.SetWindowTitle("Vous voulez rejouer ? (Appuyer sur entrée)")
	var texte string
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if g.clientInQueue != 2 {
			texte = "You are ready !\n  Waiting players " + fmt.Sprint(2-g.clientInQueue)
		} else {
			texte = "Press SPACE to restart !\n   Waiting players " + fmt.Sprint(2-g.clientInQueue)
		}
	}
	text.Draw(screen, texte, smallFont, 300, 500, globalTextColor)
	text.Draw(screen, message, smallFont, 300, 400, globalTextColor)
}

// Affichage de la grille de puissance 4, incluant les pions déjà joués.
func (g Game) drawGrid(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, globalTileSize, globalTileSize*globalNumTilesX, globalTileSize*globalNumTilesY, globalGridColor, true)

	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {

			var tileColor color.Color
			switch g.grid[x][y] {
			case p1Token:
				tileColor = globalTokenColors[g.p1Color]
			case p2Token:
				tileColor = globalTokenColors[g.p2Color]
			default:
				tileColor = globalBackgroundColor
			}

			vector.DrawFilledCircle(screen, float32(globalTileSize/2+x*globalTileSize), float32(globalTileSize+globalTileSize/2+y*globalTileSize), globalTileSize/2-globalCircleMargin, tileColor, true)
		}
	}
}
