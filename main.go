package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
    mag = 0.5
    particleSize = 4
    screenWidth = 1000
    screenHeight = 2000
    rotationRate = 3000
    gravity = 0.8*float32(particleSize)  //y acceleration (positive down)
)

var (
    vXInit        float32
    vYInit        []float32
    canvas        *ebiten.Image
    particleImage *ebiten.Image
    positions     []pos
)

// Initialise initial particle velocities
func init() {
    vXInit = 1*float32(particleSize)  //initial x velocity (positive right)
    vYInit = []float32{-5*float32(particleSize), -10.5*float32(particleSize)}  //initial y velocities (positive down)
}

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
    canvas = ebiten.NewImage(screenWidth, screenHeight)
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
    falling bool
}

func (p *particle) update() error {
    //if !p.falling {
    //    return nil
    //}
    p.vy = p.vy + gravity
    xNext := p.x + p.vx
    yNext := p.y + p.vy
    pos, falling := getParticlePosition(p, xNext, yNext)
    p.falling = falling
    //if !p.falling {
    //    p.falling = falling
    //    p.vy = 0
    //    p.vx = 0
    //}
    p.x = float32(pos.x)
    p.y = float32(pos.y)
    positions[p.id] = pos
    return nil
}

func (p *particle) draw() {
    op := &colorm.DrawImageOptions{}
    pos := positions[p.id]
    op.GeoM.Translate(float64(pos.x), float64(pos.y))
    cm := getParticleColor(p.id)
    colorm.DrawImage(canvas, particleImage, cm, op)
}

func getParticleColor(count int) colorm.ColorM {
    var cm colorm.ColorM
    cm.Scale(1.00, 0.100, 0.100, 0.5)
    theta := 2.0 * math.Pi * float64(count % rotationRate) / float64(rotationRate)
    cm.RotateHue(theta)
    return cm
}

// Create four new particles, two on either side of the given position
func (g *Game) dropParticles(x, y float32) {
    for _, side := range []int{1, -1} {
        for i, vy := range vYInit {
            xp := x + float32(side*particleSize)
            yp := y + 2*float32(i*particleSize)
            g.particles = append(g.particles, &particle{
                x: xp,
                y: yp,
                vx: float32(side)*vXInit,
                vy: vy,
                id: g.count,
                falling: true,
            })
            g.count++
            positions = append(positions, pos{
                x: int(xp/particleSize) * particleSize + int(particleSize/2),
                y: int(yp/particleSize) * particleSize + int(particleSize/2),
            })
        }
    }
}

func getParticlePosition(p *particle, targetX, targetY float32) (position pos, falling bool) {
    falling = false
    // Filter out out of bounds particles
    if p.x < 0 || p.x > screenWidth || p.y > screenHeight {
        return pos{x: int(screenWidth/2), y: screenHeight+particleSize}, falling
    }

    // Initial position in pixels
    initPOS := pos{
        x: int(p.x/particleSize) * particleSize + int(particleSize/2),
        y: int(p.y/particleSize) * particleSize + int(particleSize/2),
    }
    // Target position in pixels
    targetPOS := pos{
        x: int(targetX/particleSize) * particleSize + int(particleSize/2),
        y: int(targetY/particleSize) * particleSize + int(particleSize/2),
    }
    // Bottom position
    bottomPOS := pos{
        x: targetPOS.x,
        y: screenHeight - int(particleSize/2),
    }

    if (targetPOS.y >= bottomPOS.y) && 
        !positionOccupied(bottomPOS) {
            (*p).vx = 0
            (*p).vy = 0
            return bottomPOS, false
    } else if (initPOS == bottomPOS) || 
        (positionOccupied(pos{x: initPOS.x, y: initPOS.y-particleSize}) && 
        (positionOccupied(pos{x: initPOS.x-particleSize, y:initPOS.y-particleSize}) || positionOccupied(pos{x: initPOS.x+particleSize, y: initPOS.y-particleSize}))) {
            p.vx = 0
            p.vy = 0
            return initPOS, false
    } else if !positionOccupied(targetPOS) && (targetPOS.y < bottomPOS.y) {
        falling = true
        return targetPOS, true
    } else {
        p.vx = 0
        p.vy = max(0, particleSize)
        falling = true
        // define height range; y1 > y2
        y1 := targetPOS.y - particleSize
        y2 := initPOS.y
        for y := y1; y > y2; y -= particleSize {
            if y > bottomPOS.y {
                continue
            }
            targetPOS.y = y
            if !positionOccupied(targetPOS) {
                return targetPOS, true
            }
        }
        source := rand.New(rand.NewSource(time.Now().UnixNano()))
        i := source.Intn(2)
        //fmt.Println(i, p.falling)
        for j := 0; j < 2; j++ {
            x := targetPOS.x + (((i+j) % 2)*2 - 1) * particleSize
            if (x < int(particleSize/2)) || (x > (screenWidth - int(particleSize/2))) {
                continue
            }
            targetPOS.x = x
            if !positionOccupied(targetPOS) {
                return targetPOS, true
            }
        }
    }
    return initPOS, false
}

func positionOccupied(targetPos pos) bool {
    return slices.Contains(positions, targetPos)
}

func positionOffscreen(p pos) bool {
    return (p.x < 0) || (p.x > screenWidth) || (p.y > screenHeight)
}

func (g *Game) Update() error {
    ch := make(chan error)
    go func(particles []*particle) error {
        defer close(ch)
        for _, p := range g.particles {
            if positionOffscreen(pos{x: int(p.x), y: int(p.y)}) {
                continue
            }
            err := p.update()
            if err != nil {
                return err
            }
        }
        return nil
    }(g.particles)

    for err := range ch {
        if err != nil {
            return err
        }
    }
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()
        if !positionOffscreen(pos{x: x, y: y}) {
            g.dropParticles(float32(x), float32(y))
        }
    }
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    ch := make(chan bool)
    canvas.Clear()
    canvas.Fill(color.Black)
    go func(particles []*particle) {
        defer close(ch)
        for _, p := range particles {
            if positionOffscreen(pos{x: int(p.x), y: int(p.y)}) {
                continue
            }
            p.draw()
        }
    }(g.particles)

    for _ = range ch {
        continue
    }

    screen.DrawImage(canvas, nil)
    tps := ebiten.ActualTPS()
    cx, cy := ebiten.CursorPosition()
    msg := fmt.Sprintf("(%d, %d))\nParticles: %d\nTPS: %.2f", cx, cy, g.count, tps)
    ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return screenWidth, screenHeight
}

func main() {
    game := &Game{}
    ebiten.SetWindowSize(int(mag*screenWidth), int(mag*screenHeight))
    ebiten.SetWindowTitle("test game")
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
