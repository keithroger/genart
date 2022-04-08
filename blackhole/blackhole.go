package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"image"
	"image/color"
	"image/draw"

	"github.com/faiface/gui/win"
	"github.com/faiface/mainthread"
)

type blackHole struct {
	particles     []*particle
	eventHoriz    float64
	width, height int
}

func newBlackHole(nParticles, width, height int) *blackHole {
	w, h := float64(width), float64(height)
	particles := make([]*particle, nParticles, nParticles)
	for i := range particles {
		if i%2 == 0 {
			particles[i] = newParticle("blue")
		} else {
			particles[i] = newParticle("pink")
		}
	}

	return &blackHole{
		particles,
		math.Min(w, h) * .16,
		width,
		height,
	}
}

func (b *blackHole) cycle(drw draw.Image) {
	for i := range b.particles {
		if b.particles[i].r-b.eventHoriz < 0.01 {
			if i%2 == 0 {
				b.particles[i] = newParticle("blue")
			} else {
				b.particles[i] = newParticle("pink")
			}
		}

		b.particles[i].move(b.width, b.height, b.eventHoriz)

		b.plot(drw, b.particles[i])
	}
}

func (b *blackHole) plot(drw draw.Image, p *particle) {
	x, y := p.getCartesianCoordinates()
	setCircle(
		x+b.width/2,
		y+b.height/2,
		p.radius,
		p.colr,
		drw,
	)
}

type particle struct {
	r, theta, t, offset float64
	radius              int
	radsPerFrame        float64
	colr                color.RGBA
}

func newParticle(colr string) *particle {
	switch colr {
	case "blue":
		return &particle{
			0.0,
			0.0,
			(rand.Float64() + 0.5) * 1 * math.Pi,
			rand.Float64() * math.Pi,
			rand.Intn(3),
			(rand.Float64() + 1.0) * math.Pi / 250.0,
			color.RGBA{0, 0, 255, 255},
		}
	default:
		return &particle{
			0.0,
			0.0,
			(rand.Float64() + 0.5) * 1 * math.Pi,
			rand.Float64()*math.Pi + math.Pi,
			rand.Intn(3),
			(rand.Float64() + 1.0) * math.Pi / 250.0,
			color.RGBA{255, 105, 200, 255},
		}
	}
}

func (p *particle) move(width, height int, eventHoriz float64) {
	w, h := float64(width), float64(height)
	alpha := math.Min(w, h) / 200
	beta := 10.0 * math.Min(w, h)
	p.r = beta*math.Exp(-(p.t-alpha)/alpha) + eventHoriz
	p.theta = p.t + p.offset
	p.t = p.t*1.00005 + p.radsPerFrame
}

func (p *particle) getCartesianCoordinates() (int, int) {
	return int(p.r * math.Cos(p.theta)), int(p.r * math.Sin(p.theta))
}

func setCircle(x0, y0, r int, c color.RGBA, drw draw.Image) {
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

func run() {
	rand.Seed(time.Now().UnixNano())
	width, height := 400, 400
	blackHole := newBlackHole(1000, width, height)
	for _, p := range blackHole.particles {
		fmt.Printf("r: %f\n", p.r)
	}

	w, err := win.New(win.Title("Black Hole"), win.Size(width, height))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			w.Draw() <- func(drw draw.Image) image.Rectangle {
				r := image.Rectangle{
					image.Point{0, 0},
					image.Point{height, width},
				}
				draw.Draw(drw, r, image.Black, image.Point{0, 0}, draw.Src)

				blackHole.cycle(drw)

				return r
			}
			//time.Sleep(time.Second / 600 )
			time.Sleep(time.Second / 60)
			//time.Sleep(time.Second / 6 )
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
