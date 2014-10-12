package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/michiwend/goplaceholder"
)

func writeImageToPng(img *image.Image, name string) {
	fso, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer fso.Close()
	png.Encode(fso, (*img))
}

func main() {

	width := flag.Int("width", 0, "image width")
	height := flag.Int("height", 0, "image height")
	text := flag.String("text", "", "optional text")

	flag.Parse()

	font := "/usr/share/fonts/TTF/DejaVuSans-Bold.ttf"
	fg := color.RGBA{150, 150, 150, 255}
	bg := color.RGBA{204, 204, 204, 255}

	img, err := goplaceholder.Placeholder(*text, font, fg, bg, *width, *height)

	if err != nil {
		log.Fatal(err)
	}

	writeImageToPng(&img, "./test.png")
}
