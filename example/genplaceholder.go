package main

import (
	"image"
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

	img, err := goplaceholder.Placeholder("Lorem ipsum", 400, 200)

	if err != nil {
		log.Println(err)
	}

	writeImageToPng(&img, "./test.png")
}
