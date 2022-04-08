#Black Hole
- have event horizon where particles speed up
- add a fade property to the drawing image
    - save the last frame image and divide brightness by half and draw next frame on top
- add bool for inColor
- if particles have different scalars then they would all have to converge on the same horizon

## blackHole struct
- particles []particle
- eventHorizon float64
- displacement *mat.VecDense
- maxSpeed float64 // radians per frame

- func cycle 
    - tranforms
    - deletes old elements and replaces inplace
    - sets points to draw.image
- func appendParticle

## particle struct
- pos *mat.VecDense
- radius int
- speed float64 // measure in radians per frame
- colr color.RGBA

- func tranform

## functions
- newParticles
