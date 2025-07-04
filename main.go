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
    mag = 1.0
    particleSize = 4
    screenWidth = 1000
    screenHeight = 1000
    rotationRate = 3000
    gravity = 0.8*float32(particleSize)  //y acceleration (positive down)
)

var (
    vXInit        float32
    vYInit        []float32
    canvas        *ebiten.Image
    particleImage *ebiten.Image
    positions     []Pos
    state         string
)

// Initialise initial particle velocities
func init() {
    state = "Waiting..."
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

type Pos struct {
    x int
    y int
}

type Game struct {
    count       int
    background  int
    particles   []*Particle
    positions   []Pos
}

type Particle struct {
    x float32
    y float32
    vx float32
    vy float32
    pos Pos
    id int
    falling bool
}

func (p *Particle) update() error {
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

func (p *Particle) draw() {
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
            if positionOccupied(Pos{x: int(xp/particleSize)*particleSize+int(particleSize/2), y: int(yp/particleSize)*particleSize+int(particleSize/2)}) {
                continue
            }
            g.particles = append(g.particles, &Particle{
                x: xp,
                y: yp,
                vx: float32(side)*vXInit,
                vy: vy,
                id: g.count,
                falling: true,
            })
            g.count++
            positions = append(positions, Pos{
                x: int(xp/particleSize) * particleSize + int(particleSize/2),
                y: int(yp/particleSize) * particleSize + int(particleSize/2),
            })
        }
    }
}

func getParticlePosition(p *Particle, targetX, targetY float32) (position Pos, falling bool) {
    falling = false
    // Filter out out of bounds particles
    if p.x < 0 || p.x > screenWidth || p.y > screenHeight {
        return Pos{x: int(screenWidth/2), y: screenHeight+particleSize}, falling
    }

    // Initial position in pixels
    initPOS := Pos{
        x: int(p.x/particleSize) * particleSize + int(particleSize/2),
        y: int(p.y/particleSize) * particleSize + int(particleSize/2),
    }
    // Target position in pixels
    targetPOS := Pos{
        x: int(targetX/particleSize) * particleSize + int(particleSize/2),
        y: int(targetY/particleSize) * particleSize + int(particleSize/2),
    }
    // Bottom position
    bottomPOS := Pos{
        x: int(targetX/particleSize) * particleSize + int(particleSize/2),
        //x: targetPOS.x,
        y: screenHeight - int(particleSize/2),
    }

    if (targetPOS.y >= bottomPOS.y) && 
        !positionOccupied(bottomPOS) {
            fmt.Println("out of bounds: to bottom")
            fmt.Println(initPOS, positionOccupied(initPOS), targetX, targetY, bottomPOS, 
                positionOccupied(Pos{x:targetPOS.x, y: targetPOS.y}),
            )
            fmt.Println(positions)
            p.vx = 0
            p.vy = 0
            return bottomPOS, false
    } else if (initPOS == bottomPOS) || 
        (positionOccupied(Pos{x: initPOS.x, y: initPOS.y-particleSize}) && 
            (positionOccupied(Pos{x: initPOS.x-particleSize, y:initPOS.y-particleSize}) || 
            positionOccupied(Pos{x: initPOS.x+particleSize, y: initPOS.y-particleSize}))) {
                fmt.Println("Stationary")
                p.vx = 0
                p.vy = 0
                return initPOS, false
    } else if !positionOccupied(targetPOS) && (targetPOS.y < bottomPOS.y) {
        fmt.Println("fall")
        falling = true
        return targetPOS, true
    } else {
        fmt.Println("Fall partial")
        p.vx = 0
        p.vy = 2*particleSize  //max(0, particleSize)
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

func positionOccupied(targetPos Pos) bool {
    return slices.Contains(positions, targetPos)
}

func positionOffscreen(p Pos) bool {
    return (p.x < 0) || (p.x > screenWidth) || (p.y > screenHeight)
}

func (g *Game) Update() error {
    for _, p := range g.particles {
        if positionOffscreen(Pos{x: int(p.x), y: int(p.y)}) {
            p.falling = false
            continue
        }
        err := p.update()
        if err != nil {
            return err
        }
    }

    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()
        if !positionOffscreen(Pos{x: x, y: y}) && !positionOccupied(Pos{x: int(x/particleSize)*particleSize+int(particleSize/2), y: int(y/particleSize)*particleSize+int(particleSize/2)}) {
            g.dropParticles(float32(x), float32(y))
        }
    }
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    canvas.Clear()
    canvas.Fill(color.Black)
    falling := false
    for _, p := range g.particles {
        if p.falling {
            falling = p.falling
        }
        if positionOffscreen(Pos{x: int(p.x), y: int(p.y)}) {
            continue
        }
        p.draw()
    }
    if falling && (state != "Falling...") {
        state = "Falling..."
        //log.Println("Falling...")
    } else if !falling && (state == "Falling...") {
        state = "Settled."
    }
    falling = false

    screen.DrawImage(canvas, nil)
    tps := ebiten.ActualTPS()
    cx, cy := ebiten.CursorPosition()
    msg := fmt.Sprintf("(%d, %d))\nParticles: %d\nTPS: %.2f\nState: %v", cx, cy, g.count, tps, state)
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
