# nes

This is an example on how to register a custom image format and reading NES ROMs. After importing it with `_` you can use it with `image.Decode`. It uses black, red, green and blue as colors and shows 16 tiles per row.

```go
import (
	"image"
	"os"

	_ "github.com/BakeRolls/nes"
)

func main() {
	r, err := os.Open("rom.nes")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatal(err)
	}
}
```

To view it in action, see the [NES viewer](/cmd/nes-viewer/main.go) which opens a window displaying the sprites of a given ROM.
