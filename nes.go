// Package nes registers a custom image format for reading NES ROMs. After importing it with `_` you can use `image.Decode` to decode a roms sprites into an `image.Image`.
//
// (image/Super Mario Bros 1 Sprites) ./testdata/super_mario_bros_sprites.png
package nes

import (
	"image"
	"image/color"
	"io"

	"github.com/pkg/errors"
)

const (
	header = "NES\x1a"

	headerSize     = 16
	spriteSize     = 16
	prgBankSize    = 16 * 1024
	chrBankSize    = 8 * 1024
	spritesPerBank = chrBankSize / spriteSize

	spriteWidth = 8

	spritesPerRow = 16
)

var colors = []color.Color{
	color.RGBA{},
	color.RGBA{R: 255, A: 255},
	color.RGBA{G: 255, A: 255},
	color.RGBA{B: 255, A: 255},
}

// Decode reads a NES ROM and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	c, prgBanks, chrBanks, err := decodeConfig(r)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))

	prg := make([]byte, prgBanks*prgBankSize)
	if _, err := io.ReadFull(r, prg); err != nil {
		return nil, errors.Wrap(err, "unexpected EOF in PRG ROM")
	}

	for i := 0; i < chrBanks; i++ {
		for j := 0; j < spritesPerBank; j++ {
			s := make([]byte, spriteSize)
			if _, err := io.ReadFull(r, s); err != nil {
				if err == io.EOF {
					return nil, errors.Wrap(err, "unexpected EOF in sprite")
				}
				return nil, errors.Wrapf(err, "could not read sprite %d", j)
			}

			o := i*spritesPerBank + j
			for k, c := range decode(s) {
				x := (o % spritesPerRow * spriteWidth) + (k % spriteWidth)
				y := (o / spritesPerRow * spriteWidth) + (k / spriteWidth)
				img.Set(x, y, colors[c])
			}
		}
	}

	return img, nil
}

func decode(s []byte) []byte {
	cs := make([]byte, spriteWidth*spriteWidth)
	for i := uint(0); i < spriteWidth; i++ {
		c1 := s[i]
		c2 := s[i+spriteWidth]
		for j := uint(0); j < spriteWidth; j++ {
			v := (c1 >> j & 1) + (c2>>j)&1<<1
			x := spriteWidth - j - 1
			y := i * spriteWidth
			cs[x+y] = v
		}
	}
	return cs
}

func decodeConfig(r io.Reader) (image.Config, int, int, error) {
	h := make([]byte, headerSize)
	if _, err := io.ReadFull(r, h); err != nil {
		return image.Config{}, 0, 0, errors.Wrap(err, "could not read NES header")
	}
	if h[5] == 0 {
		return image.Config{}, 0, 0, errors.New("no tiles in CHR ROM")
	}
	config := image.Config{
		ColorModel: color.RGBAModel,
		Width:      spritesPerRow * spriteWidth,
		Height:     int(h[5]) * spritesPerBank / spritesPerRow * spriteWidth,
	}
	return config, int(h[4]), int(h[5]), nil
}

// DecodeConfig returns the color model and dimensions of a NES ROM without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	c, _, _, err := decodeConfig(r)
	return c, err
}

func init() {
	image.RegisterFormat("nes", header, Decode, DecodeConfig)
}
