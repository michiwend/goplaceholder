package main

import (
	"image/color"
	"image/png"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/michiwend/goplaceholder"
)

func serveImage(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	r.ParseForm()
	text := r.FormValue("text")
	width, _ := strconv.ParseInt(params["width"], 10, 32)
	height, _ := strconv.ParseInt(params["height"], 10, 32)
	// FIXME err handling

	if width > 4000 || height > 4000 {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("image too large"))
		return
	}

	img, err := goplaceholder.Placeholder(
		text,
		"/usr/share/fonts/TTF/DejaVuSans-Bold.ttf",
		color.RGBA{150, 150, 150, 255},
		color.RGBA{204, 204, 204, 255},
		int(width), int(height))

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)

	imgName := params["width"]
	if w, ok := params["height"]; ok {
		imgName += "x" + w
	}
	imgName += ".png"

	log.WithFields(log.Fields{
		"width":      width,
		"height":     height,
		"foreground": fg,
		"background": bg,
		"text":       text,
	}).Infof("Served placeholder image \"%s\"", imgName)

}

func main() {

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{width:[0-9]+}x{height:[0-9]+}.png", serveImage).Methods("GET")
	rtr.HandleFunc("/{width:[0-9]+}.png", serveImage).Methods("GET")

	http.Handle("/", rtr)

	log.Info("Starting placeholder service on port 3000...")
	log.Fatal(http.ListenAndServe("127.0.0.1:3000", nil))

}
