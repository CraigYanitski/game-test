package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const (
    particleSize = 4
    screenWidth = 320
    screenHeight = 320
    vXInit = float32(0.25)  //initial x velocity (positive right)
    vYInit = []float32{-0.25, -0.50}  //initial y velocities (positive down)
    gravity = 0.15  //y acceleration (positive down)
)

var (
    canvas        *ebiten.Image
    particleImage *ebiten.Image
    positions     []pos
)

// initialise particle image
func init() {
    // Define three possible alpha values
    const (
        a0 = 0x40
        a1 = 0xc0
        a2 = 0xff
    )
    // "Draw" the alpha map for each particle
    pixels := []uint8{
        a0, a1, a1, a0,
        a1, a2, a2, a1,
        a1, a2, a2, a1,
        a0, a1, a1, a0,
    }
    // Use alpha map in particle image
    particleImage = ebiten.NewImageFromImage(&image.Alpha{
        Pix: pixels,
        Stride: 4,
        Rect: image.Rect(0, 0, 4, 4),
    })
}

// initialise canvas image
func init() {
    canvas := &ebiten.NewImage(screenWidth, screenHeight)
}

type pos struct {
    x int
    y int
}

type Game struct {
    count       int
    background  int
    particles   []*particle
    //positions   []pos
}

type particle struct {
    x float32
    y float32
    vx float32
    vy float32
    pos pos
    id int
}

func (p *particle) update() error {
    // TODO: finish particle update
    p.vy := p.vy + gravity
    xNext := p.x + p.vx
    yNext := p.y + p.vy
    pos := getParticlePosition(p.x, p.y, xNext, yNext)
    return nil
}

func (p *particle) draw() {
    op := &colorm.DrawImageOptions{}
    pos = positions[p.id]
    op.GeoM.Translate(float64(pos.x), float64(pos.y))
    cm := getParticleColor(p.id)
    colorm.DrawImage(canvas, particleImage, cm, op)
}

func getParticleColor(count int) colorm.ColorM {
    var cm colorm.ColorM
    cm.Scale(0.5, 1.0, 1.0, 1.0)
    tps := ebiten.TPS()
    theta := 2.0 * math.Pi * float32(count%tps) / float32(tps)
    cm.RotateHue(theta)
    return cm
}

// Create four new particles, two on either side of the given position
func (g *Game) dropParticles(x, y float32) {
    for _, side := range []int{1, -1} {
        for i, vy := range vYInit {
            g.particles = append(g.particles, &particle{
                x: x + float32(side*particleSize),
                y: y + float32(i*particleSize),
                vx: float32(side)*vXInit,
                vy: vy,
                id: g.count,
            })
            g.count++
        }
    }
}

func getParticlePosition(initX, initY, targetX, targetY float32) (pos) {
    // TODO: finish function to return the particle position
    initPOS := pos{
        x: int(initX/particleSize) * particleSize + int(particleSize/2)
        y: int(initY/particleSize) * particleSize + int(particleSize/2)
    }
    targetPOS := pos{
        x: int(targetX/particleSize) * particleSize + int(particleSize/2),
        y: int(targetY/particleSize) * particleSize + int(particleSize/2),
    }

    if !positionOccupied(targetPOS) {
        return targetPOS
    } else if !positionOccupied(initPOS) {
    } else {
    }
    return pos{x: 0, y: 0}
}

func (g *Game) positionOccupied(targetPos pos) bool {
    // Function to determine if particle is already in a given position
    //targetPos := []int{x, y}
    return slices.Contains(g.positions, targetPos)
    
    //for _, p := range g.positions {
    //    if slices.Equal(p, targetPos) {
    //        return true
    //    }
    //}
    //return false
}

func (g *Game) Update() error {
    // TODO: extend game update loop
    for _, p := range g.particles {
        err := p.update()
        if err != nil {
            return err
        }
    }
    if ebiten.IsMouseButtonPressed() {
        x, y := ebiten.CursorPosition()
        g.dropParticles(float32(x), float32(y))
    }
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // TODO: Implement game drawing logic
    canvas.Clear()
    canvas.Fill(color.Black)
    for _, particle := range g.particles {
        particle.draw()
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return screenWidth, screenHeight
}

func main() {
    game := &Game{}
    ebiten.SetWindowSize(2+screenWidth, 2*screenHeight)
    ebiten.SetWindowTitle("test game")
    if err := ebiten.RunGame(game); if err != nil {
        log.Fatal(err)
    }
}
