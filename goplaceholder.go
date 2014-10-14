/*
 * Copyright (c) 2014 Michael Wendland
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 *	Authors:
 *		Michael Wendland <michael@michiwend.com>
 */

/*
package goplaceholder implements a simple library to generate placeholder
images using freetype-go.
*/
package goplaceholder

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"strconv"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/raster"
)

const (
	maxTextBoundsToImageRatioY = 0.23
	maxTextBoundsToImageRatioX = 0.64
	dpi                        = 72.00
	testFontSize               = 1.00 // use heigher values (>=100) when hinting enabled
)

// Placeholder returns a placeholder image with the given text or, if text was
// an empty string, with the image bounds in the form "800x600".
func Placeholder(text, ttfPath string, foreground, background color.RGBA, width, height int) (image.Image, error) {

	if width < 0 || height < 0 {
		return nil, errors.New("negative values not allowed")
	}
	if width == 0 && height == 0 {
		return nil, errors.New("either width or height needs to be > 0")
	}

	if width == 0 {
		width = height
	} else if height == 0 {
		height = width
	}

	if text == "" {
		text = strconv.Itoa(width) + " x " + strconv.Itoa(height)
	}

	fontBytes, err := ioutil.ReadFile(ttfPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	fg_img := image.NewUniform(foreground)
	bg_img := image.NewUniform(background)

	testImg := image.NewRGBA(image.Rect(0, 0, 0, 0))

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(testFontSize)
	c.SetSrc(fg_img)
	c.SetDst(testImg)
	c.SetClip(testImg.Bounds())
	c.SetHinting(freetype.NoHinting)

	// first draw with testFontSize to get the text extent
	var textExtent raster.Point
	drawPoint := freetype.Pt(0, int(c.PointToFix32(testFontSize)>>8))
	textExtent, err = c.DrawString(text, drawPoint)
	if err != nil {
		return nil, err
	}

	// calculate font scales to stay within the bounds
	scaleX := float64(c.PointToFix32(float64(width)*maxTextBoundsToImageRatioX)) / float64(textExtent.X)
	scaleY := float64(c.PointToFix32(float64(height)*maxTextBoundsToImageRatioY)) / float64(textExtent.Y)

	fontsize := testFontSize * math.Min(scaleX, scaleY)

	// draw with scaled fontsize to get the real text extent. This could also be
	// done by scaling up the textExtent from the previous drawing but it's less
	// precise.
	c.SetFontSize(fontsize)
	drawPoint = freetype.Pt(0, 0)
	textExtent, err = c.DrawString(text, drawPoint)
	if err != nil {
		return nil, err
	}

	// finally draw the centered text
	drawPoint = freetype.Pt(
		int(c.PointToFix32(float64(width)/2.0)-textExtent.X/2)>>8,
		int(c.PointToFix32(float64(height)/2.0+fontsize/2.6))>>8)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), bg_img, image.ZP, draw.Src)

	c.SetDst(img)
	c.SetClip(img.Bounds())
	_, err = c.DrawString(text, drawPoint)
	if err != nil {
		return nil, err
	}

	return img, nil
}
