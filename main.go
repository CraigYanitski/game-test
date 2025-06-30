package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
    particleSize = 4
    screenWidth = 500
    screenHeight = 500
    gravity = 1.1*float32(particleSize)  //y acceleration (positive down)
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
}

func (p *particle) update() error {
    // TODO: finish particle update
    if p.vx == 0 {
        return nil
    }
    p.vy = p.vy + gravity
    xNext := p.x + p.vx
    yNext := p.y + p.vy
    pos, falling := getParticlePosition(p.x, p.y, xNext, yNext)
    if !falling {
        p.vy = 0
        p.vx = 0
    }
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
    cm.Scale(0.75, 0.125, 0.125, 0.75)
    tps := ebiten.TPS()
    theta := 1.0 * math.Pi * float64(count%tps) / float64(tps)
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
            })
            g.count++
            positions = append(positions, pos{
                x: int(xp/particleSize) * particleSize + int(particleSize/2),
                y: int(yp/particleSize) * particleSize + int(particleSize/2),
            })
        }
    }
}

func getParticlePosition(initX, initY, targetX, targetY float32) (position pos, falling bool) {
    falling = false
    // Filter out out of bounds particles
    if initX < 0 || initX > screenWidth || initY > screenHeight {
        return pos{x: int(screenWidth/2), y: screenHeight+particleSize}, falling
    }

    // Initial position in pixels
    initPOS := pos{
        x: int(initX/particleSize) * particleSize + int(particleSize/2),
        y: int(initY/particleSize) * particleSize + int(particleSize/2),
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

    if (targetPOS.y >= (screenHeight - int(particleSize/2))) && 
        !positionOccupied(bottomPOS) {
            return bottomPOS, falling
    } else if !positionOccupied(targetPOS) && (targetPOS.y < (screenHeight - int(particleSize/2))) {
        falling = true
        return targetPOS, falling
    //} else if positionOccupied(initPOS) {
    //    // define height range; y1 < y2
    //    y1 := initPOS.y - screenHeight
    //    y2 := int(particleSize/2)
    //    i := rand.IntN(2)
    //    xTarget := targetPOS.x
    //    for y := y1; y >= y2; y -= particleSize {
    //        // Check position directly above previous target
    //        targetPOS.y = y
    //        if !positionOccupied(targetPOS) {
    //            return targetPOS
    //        }
    //        // Otherwise check adjacent positions
    //        for j := 0; j < 2; j++ {
    //            x := []int{targetPOS.x+particleSize, targetPOS.x-particleSize}[(i+j)%2]
    //            targetPOS.x = x
    //            if !positionOccupied(targetPOS) {
    //                return targetPOS
    //            }
    //        }
    //        // Reset target x position
    //        targetPOS.x = xTarget
    //    }
    } else {
        falling = true
        // define height range; y1 > y2
        y1 := targetPOS.y - screenHeight
        y2 := initPOS.y
        for y := y1; y > y2; y -= screenHeight {
            if y > (screenHeight - int(particleSize/2)) {
                continue
            }
            targetPOS.y = y
            if !positionOccupied(targetPOS) {
                return targetPOS, falling
            }
        }
        i := rand.IntN(2)
        for j := 0; j < 2; j++ {
            x := []int{targetPOS.x+particleSize, targetPOS.x-particleSize}[(i+j)%2]
            if (x < int(particleSize/2)) || (x > (screenWidth - int(particleSize/2))) {
                continue
            }
            targetPOS.x = x
            if !positionOccupied(targetPOS) {
                return targetPOS, falling
            }
        }
    }
    return initPOS, falling
}

func positionOccupied(targetPos pos) bool {
    // Function to determine if particle is already in a given position
    //targetPos := []int{x, y}
    return slices.Contains(positions, targetPos)
    
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
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()
        g.dropParticles(float32(x), float32(y))
    }
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    canvas.Clear()
    canvas.Fill(color.Black)
    for _, particle := range g.particles {
        if particle.y > (float32(screenHeight) - particleSize/2) {
            continue
        }
        particle.draw()
    }

    screen.DrawImage(canvas, nil)
    tps := ebiten.TPS()
    cx, cy := ebiten.CursorPosition()
    msg := fmt.Sprintf("(%d, %d))\nParticles: %d\nTPS: %v", cx, cy, g.count, tps)
    ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return screenWidth, screenHeight
}

func main() {
    game := &Game{}
    ebiten.SetWindowSize(2*screenWidth, 2*screenHeight)
    ebiten.SetWindowTitle("test game")
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
