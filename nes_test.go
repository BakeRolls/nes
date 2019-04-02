package nes_test

import (
	"fmt"
	"image"
	"log"
	"os"

	_ "github.com/bakerolls/nes"
)

func ExampleDecode() {
	// This rom file does not contain any game logic.
	r, err := os.Open("testdata/super_mario_bros_sprites.nes")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	_, format, err := image.Decode(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(format)
	// Output: nes
}
