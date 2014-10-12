package goplaceholder

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"strconv"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/raster"
)

const (
	maxTextBoundsToImageRatioY = 0.23
	maxTextBoundsToImageRatioX = 0.44
	dpi                        = 72.00
	testFontSize               = 1.00 // use heigher values (>=100) when hinting enabled
)

// FIXME foreground / background and font params
func Placeholder(text string, width, height int) (image.Image, error) {

	if text == "" {
		text = strconv.Itoa(width) + " x " + strconv.Itoa(height)
	}

	ttfPath := "/usr/share/fonts/TTF/DejaVuSans-Bold.ttf"
	//ttfPath := "/home/michael/.local/share/fonts/DejaVu Sans Mono Bold for Powerline.ttf"

	foreground := image.NewUniform(color.RGBA{150, 150, 150, 255})
	background := image.NewUniform(color.RGBA{204, 204, 204, 255})

	fontBytes, err := ioutil.ReadFile(ttfPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	testImg := image.NewRGBA(image.Rect(0, 0, 0, 0))

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(testFontSize)
	c.SetSrc(foreground)
	c.SetDst(testImg)
	c.SetClip(testImg.Bounds())
	c.SetHinting(freetype.NoHinting)

	var textExtent raster.Point
	drawPoint := freetype.Pt(0, int(c.PointToFix32(testFontSize)>>8))
	textExtent, err = c.DrawString(text, drawPoint)
	if err != nil {
		return nil, err
	}

	scaleX := float64(c.PointToFix32(float64(width)*maxTextBoundsToImageRatioX)) / float64(textExtent.X)
	scaleY := float64(c.PointToFix32(float64(height)*maxTextBoundsToImageRatioY)) / float64(textExtent.Y)

	//fontsize := testFontSize * math.Min(scaleX, scaleY)

	var fontsize, originX, originY float64

	if scaleX < scaleY {
		fmt.Println("extending to X-bounds")
		fontsize = testFontSize * scaleX
		originX = (float64(width) - float64(width)*maxTextBoundsToImageRatioX) / 2.0
	} else {
		fmt.Println("extending to Y-bounds")
		fontsize = testFontSize * scaleY
		// FIXME
		originX = float64(width)/2.0 - (float64(textExtent.X>>8)*fontsize)/2.0
	}
	originY = float64(height)/2.0 + fontsize/2.0

	drawPoint = freetype.Pt(
		int(c.PointToFix32((originX))>>8),
		int(c.PointToFix32((originY))>>8))

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), background, image.ZP, draw.Src)

	c.SetDst(img)
	c.SetClip(img.Bounds())
	c.SetFontSize(fontsize)
	_, err = c.DrawString(text, drawPoint)
	if err != nil {
		return nil, err
	}

	return img, nil
}
