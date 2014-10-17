package main

import (
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/michiwend/goplaceholder"
)

func serveImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")

	params := mux.Vars(r)
	r.ParseForm()
	text := r.FormValue("text")
	width, _ := strconv.ParseInt(params["width"], 10, 32)
	height, _ := strconv.ParseInt(params["height"], 10, 32)
	// FIXME err handling

	img, _ := goplaceholder.Placeholder(
		text,
		"/usr/share/fonts/TTF/DejaVuSans-Bold.ttf",
		color.RGBA{150, 150, 150, 255},
		color.RGBA{204, 204, 204, 255},
		int(width), int(height))

	png.Encode(w, img)
	log.Printf("served image: w: %d h: %d\n", width, height)
}

func main() {

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{width:[0-9]+}x{height:[0-9]+}.png", serveImage).Methods("GET")
	rtr.HandleFunc("/{width:[0-9]+}.png", serveImage).Methods("GET")

	http.Handle("/", rtr)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":3000", nil))

}
