package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"time"
    "math/rand"

	"github.com/faiface/gui/win"
	"github.com/faiface/mainthread"

	"github.com/keithroger/imgge"
	"github.com/keithroger/txtplacer"
)

const (
	screenW, screenH = 1080, 1920
	// screenW, screenH = 800, 800
	maxEffectFrames  =  5
    fontfile = "playfair.ttf"
	poem             = "Twinkle, twinkle, little bat!\n" +
		"How I wonder what you're at!\n" +
		"Up above the world you fly,\n" +
		"Like a tea-tray in the sky.\n" +
		"Twinkle, twinkle, little bat!\n" +
		"How I wonder what you're at!`"
)

func run() {
	img := image.NewRGBA(image.Rect(0, 0, screenW, screenH))

	// Create effects
    effects := []struct{
        effect imgge.Effect
        frames int
    }{
        {imgge.NewShift(img.Bounds(), 15, 2, 100), 0},
        {imgge.NewColorShift(img.Bounds(), 10, 5, 100), 0},
        {imgge.NewPixelSort(img.Bounds(), 40, 200, "horiz"), 0},
        {imgge.NewPixelPop(img.Bounds(), 5, 10, 500), 0},
    }

	// Create txtplacer
	placer, err := txtplacer.NewPlacer(img, fontfile, 48.0)
	placer.SetColor(color.White)
	if err != nil {
		log.Fatalf("Font file failed to load: %v", err)
	}

	// Create window
	w, err := win.New(win.Title("Bat"), win.Size(screenW, screenH))
	if err != nil {
		panic(err)
	}

	go func() {
		for i := 0; ; i++ {
			draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)
			placer.WriteAtCenter(poem, 1000, "left")


            for i := range effects {
                if effects[i].frames == 0 {
                    effects[i].frames = rand.Intn(maxEffectFrames) + 1
                    effects[i].effect.Randomize()
                }
                effects[i].frames--

                effects[i].effect.Apply(img)

            }

			w.Draw() <- func(drw draw.Image) image.Rectangle {
				draw.Draw(drw, drw.Bounds(), img, image.Point{}, draw.Src)
				return drw.Bounds()
			}

			time.Sleep(time.Second / 100)
		}
	}()

	for event := range w.Events() {
		switch event.(type) {
		case win.WiClose:
			close(w.Draw())
		}
	}
}

func main() {
	mainthread.Run(run)
}
