package main

import (
	"errors"
	"image/color"
	"image/png"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/michiwend/goplaceholder"
)

func hexToRGB(h string) (uint8, uint8, uint8, error) {
	rgb, err := strconv.ParseUint(string(h), 16, 32)
	if err == nil {
		return uint8(rgb >> 16), uint8((rgb >> 8) & 0xFF), uint8(rgb & 0xFF), nil
	}
	return 0, 0, 0, err
}

func normalizeHex(h string) string {
	h = strings.TrimPrefix(h, "#")
	if len(h) != 3 && len(h) != 6 {
		return ""
	}
	if len(h) == 3 {
		h = h[:1] + h[:1] + h[1:2] + h[1:2] + h[2:] + h[2:]
	}
	return h
}

func paramToColor(param, defaultValue string) (color.RGBA, error) {

	if len(param) == 0 {
		param = defaultValue
	}

	hexColor := normalizeHex(param)
	if len(hexColor) == 0 {
		return color.RGBA{}, errors.New("bad hex color format")
	}

	R, G, B, err := hexToRGB(hexColor)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{R, G, B, 255}, nil
}

func serveImage(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	r.ParseForm()

	width, _ := strconv.ParseInt(params["width"], 10, 32)
	height, _ := strconv.ParseInt(params["height"], 10, 32)
	// FIXME err handling

	text := r.FormValue("text")

	if width > 4000 || height > 4000 {
		http.Error(w, "Image too large", http.StatusRequestEntityTooLarge)
		log.WithFields(log.Fields{
			"width":  width,
			"height": height,
		}).Warn("requested image too large")
		return
	}

	foregroundValue := r.FormValue("fg")
	backgroundValue := r.FormValue("bg")

	fg, err := paramToColor(foregroundValue, "969696")
	if err != nil {
		http.Error(w, "Bad value for foreground color", http.StatusBadRequest)
		log.WithField("color", foregroundValue).Error(err)
		return
	}
	bg, err := paramToColor(backgroundValue, "CCCCCC")
	if err != nil {
		http.Error(w, "Bad value for background color", http.StatusBadRequest)
		log.WithField("color", backgroundValue).Error(err)
		return
	}

	img, err := goplaceholder.Placeholder(
		text,
		"/usr/share/fonts/TTF/DejaVuSans-Bold.ttf",
		fg, bg,
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
