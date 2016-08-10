package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
)

var frames []image.Image

const perRow = 8

var outFile = flag.String("out", "out.png", "png file to write sprite sheet")

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("SPRITE STITCHER\nUsage: spritesheet --out file.png in1.{gif,png} in2.png ...")
		os.Exit(1)
	}

	for _, name := range flag.Args() {
		frames = append(frames, decode(name)...)
	}

	w, h := spriteBounds()
	fmt.Println("WIDTH", w, "HEIGHT", h)
	rows := len(frames) / perRow

	canvas := image.NewRGBA(image.Rect(0, 0, perRow*w, h*rows))

	x := 0
	y := 0
	for _, frame := range frames {
		dp := image.Point{x * w, y * h}
		fb := frame.Bounds()

		r := image.Rectangle{dp.Add(fb.Min), dp.Add(fb.Max)}
		draw.Draw(canvas, r, frame, fb.Min, draw.Over)

		x++
		if x >= perRow {
			y++
			x = 0
		}
	}

	out, err := os.Create(*outFile)
	panicerr(err)
	panicerr(png.Encode(out, canvas))
}

func spriteBounds() (w, h int) {
	bounds := frames[0].Bounds()
	w, h = bounds.Size().X, bounds.Size().Y
	for _, f := range frames {
		bounds = f.Bounds()
		if bounds.Size().X > w {
			w = bounds.Size().X
		}
		if bounds.Size().Y > h {
			h = bounds.Size().Y
		}
	}
	return
}

func decode(name string) []image.Image {
	ext := filepath.Ext(name)
	f, err := os.Open(name)
	panicerr(err)
	defer f.Close()

	switch ext {
	case ".gif":
		img, err := gif.DecodeAll(f)
		panicerr(err)
		imgs := make([]image.Image, 0, len(img.Image))
		for _, sub := range img.Image {
			imgs = append(imgs, sub)
		}
		return imgs
	case ".png":
		img, err := png.Decode(f)
		panicerr(err)
		return []image.Image{img}
	}
	panic("bad ext: " + ext)
}

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}
