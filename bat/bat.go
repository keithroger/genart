package main

import (
    "fmt"
	"image"
	"image/draw"
	"log"
	"os"
    "image/color"
    "time"

	"github.com/faiface/gui/win"
	"github.com/faiface/mainthread"

	"github.com/keithroger/imgge/effects"
	"github.com/keithroger/txtplacer"
)

const (
    screenW, screenH = 1080, 1920
    // screenW, screenH = 800, 800
)
func run() {
    img := image.NewRGBA(image.Rect(0, 0, screenW, screenH))

    // effects
    shift := effects.NewShift(img, 15, 5, 200)
    cshift := effects.NewColorShift(img, 10, 10, 200)
    pixelSort := effects.NewPixelSort(img, 40, 2000, "horiz")
    pixelPop := effects.NewPixelPop(img, 3, 10, 500)

    // Create txtplacer
    placer, err := txtplacer.NewPlacer(img, "playfair.ttf", 48.0)
    placer.SetColor(color.White)
    if err != nil {
        log.Fatalf("Font file failed to load: %v", err)
    }

    // Read txt file
    inFile, err := os.ReadFile("poem.txt")
    if err != nil {
        log.Fatalf("Failed to read poem file: %v", err)
    }

    w, err := win.New(win.Title("Bat"), win.Size(screenW, screenH))
    if err != nil {
        panic(err)
    }

    text := string(inFile)
    fmt.Println(text)

    go func() {
        for i := 0; ; i++{
            draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)
            placer.WriteAtCenter(text, 1000, "center")

            // effects
            // randomize effects after amout of time
            shift.Apply(img)
            cshift.Apply(img)
            pixelSort.Apply(img)
            pixelPop.Apply(img)

            w.Draw() <- func(drw draw.Image) image.Rectangle {
                draw.Draw(drw, drw.Bounds(), img, image.Point{}, draw.Src)
                return drw.Bounds()
            }

            if i > 1000 {
                i = 0
            }
            time.Sleep(time.Second/100)
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
