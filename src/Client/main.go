package main

import (
	"log"
	"net"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font/opentype"
)

// Mise en place des polices d'écritures utilisées pour l'affichage.
func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 30,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	largeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 50,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Création d'une image annexe pour l'affichage des résultats.
func init() {
	offScreenImage = ebiten.NewImage(globalWidth, globalHeight)
}

// Création, paramétrage et lancement du jeu.
func main() {
	//connection
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Dial error : ", err)
		return
	}
	defer conn.Close()

	log.Println("Je suis connecté")
	//connection
	g := game{}
	
	ebiten.SetWindowTitle("Programmation système : projet puissance 4")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
// penser au routine
// etablir connection entre les client
// ensuite un client joue avec un autre client