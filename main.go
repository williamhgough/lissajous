package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

const (
	whiteIndex = 0
	blackIndex = 1
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lissajous(w)
	})
	http.ListenAndServe(":8080", nil)
}

func lissajous(out io.Writer) {
	happy := colorful.FastHappyColor()

	palette := []color.Color{
		color.White,
		happy,
	}

	const (
		// Number of complete X oscillator revolutions
		cycles = 5
		// Angular resolution
		res = 0.001
		// Image canvas covers [-size.. +size]
		size = 100
		// Number of animation frames
		nframes = 128
		// Delay between frames in 10ms units
		delay = 4
	)

	// relative frequency of Y oscillator
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}

	// phase difference
	phase := 0.0

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)

		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
