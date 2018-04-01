// Open a window and display the tiles from the passed NES ROM.
// By importing image/png this could be used as a starting point for an very
// simple image viewer.

package main

import (
	"image"
	"log"
	"os"

	"github.com/BakeRolls/canvas"
	_ "github.com/BakeRolls/nes"
)

func main() {
	r, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatal(err)
	}
	c, err := canvas.New(img, 2, os.Args[1])
	for c.Update() {
		c.Draw()
	}
}
