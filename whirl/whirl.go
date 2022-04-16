package main

import (
    "math"
    "math/rand"
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "time"

	"github.com/faiface/gui/win"
	"github.com/faiface/mainthread"
)

const (
    screenW, screenH = 800, 800
    graphW, graphH = 50.0, 50.0
    nParticles = 500
    speed = 0.025
    fade = 200 
    maxSteps = 25
    padW, padH = screenW/10, screenH/10 
    //theta = 0.001
)

func NewParticle() particle {
    x, y := graphW * (rand.Float64() - 0.5), graphH * (rand.Float64() - 0.5)

    return particle{x, y, rand.Intn(50)}
}

type particle struct {
    x, y float64
    steps int
}

func (p *particle) next() {
    p.x, p.y = p.x + dx(p.x, p.y), p.y + dy(p.x, p.y)
    p.steps++
    if p.steps > maxSteps {
        *p = NewParticle()
    }
}

func drawParticle(img draw.Image, p *particle) {
    p.next()

    imgX := int((p.x + 0.5*graphW)*(float64(screenW + 2.0*padW)/graphW))
    imgY := int((p.y + 0.5*graphH)*(float64(screenH + 2.0*padH)/graphH))

    
    img.Set(imgX, imgY, color.White) 
    // setCircle(imgX, imgY, 3, color.White, img)
}

func dx(x, y float64) float64 {
    return speed * math.Exp(-1.0/norm(x,y)) * (y - x)
}

func dy(x, y float64) float64 {
    return speed * math.Exp(-1.0/norm(x,y)) * (-y - x)
}

func norm(x, y float64) float64{
    return math.Sqrt(x*x + y*y)
}

func setCircle(x0, y0, r int, c color.Color, drw draw.Image) {
	x, y := r, 0

	if r > 0 {
		drw.Set(x+x0, -y+y0, c)
		drw.Set(y+x0, x+y0, c)
		drw.Set(-y+x0, x+y0, c)
	}
	P := 1 - r
	for x > y {
		y++

		if P <= 0 {
			P = P + 2*y + 1
		} else {
			x--
			P = P + 2*y - 2*x + 1
		}

		if x < y {
			break
		}

		drw.Set(x+x0, y+y0, c)
		drw.Set(-x+x0, y+y0, c)
		drw.Set(x+x0, -y+y0, c)
		drw.Set(-x+x0, -y+y0, c)

		if x != y {
			drw.Set(y+x0, x+y0, c)
			drw.Set(-y+x0, x+y0, c)
			drw.Set(y+x0, -x+y0, c)
			drw.Set(-y+x0, -x+y0, c)
		}
	}
}

// TODO add particles from different areas
func run() {
    rand.Seed(time.Now().UnixNano())

    fmt.Println("running..")

    particles := make([]particle, nParticles)
    for i := range particles {
        particles[i] = NewParticle()
    }

    w, err := win.New(win.Title("Whirl"), win.Size(screenW, screenH))
    if err != nil {
        panic(err)
    }

    img := image.NewGray(image.Rect(-padW, -padH, screenW + padW, screenH + padH))
    alpha := &image.Uniform{color.Alpha{255 / fade}}
    // draw.Draw(alpha, alpha.Bounds(), &image.Uniform{color.Gray{2}}, image.Point{}, draw.Src)

    go func() {
        for {
            // fading with alpha overlay mask
            draw.DrawMask(img, img.Bounds(), image.Black, image.Point{}, alpha, image.Point{}, draw.Over)
            for i := range particles {
                drawParticle(img, &particles[i])
            }

            w.Draw() <- func(drw draw.Image) image.Rectangle {

                draw.Draw(drw, drw.Bounds(), img, image.Point{padW, padH}, draw.Src)


                return drw.Bounds()
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
