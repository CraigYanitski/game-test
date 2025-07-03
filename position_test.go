package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParticleMovement(t *testing.T) {
    //const (
    //    particleSize = 4
    //    gravity = 4
    //)
    positions = []Pos{
        Pos{x: 2, y: 998},
        Pos{x: 2, y: 994},
        Pos{x: 2, y: 990},
        Pos{x: 6, y: 998},
        Pos{x: 6, y: 994},
        Pos{x: 6, y: 990},
        Pos{x: 10, y: 998},
        Pos{x: 10, y: 990},
        Pos{x: 14, y: 998},
        Pos{x: 14, y: 994},
        Pos{x: 14, y: 990},
        Pos{x: 18, y: 998},
        Pos{x: 18, y: 994},
    }
    particles := make([]*Particle, len(positions))
    falling := false
    for i, pos := range positions {
        p := &Particle{
            x: float32(pos.x),
            y: float32(pos.y),
            vx: 0.0,
            vy: 0.0,
        }
        particles[i] = p
        //fmt.Println(*particles[i])
    }
    fmt.Println(particles)
    for i, p := range particles {
        fmt.Printf("Free? %v\n", positionOccupied(Pos{x: int(p.x), y: int(p.y)}))
        pos, f := getParticlePosition(p, p.x+p.vx, p.y+p.vy+float32(particleSize))
        fmt.Printf("particle position %v, %v -> %v, %v\n", int(p.x), int(p.y), pos.x, pos.y)
        p.x = float32(pos.x)
        p.y = float32(pos.y)
        positions[i] = pos
        if f {
            falling = f
        }
    }

    if falling {
        fmt.Println("pixel moved")
    }
    fmt.Println(screenHeight)
    
    require.NotNil(t, particles[4])
    assert.Equal(t, positions[0], Pos{x: 2, y: 998})
    assert.Equal(t, positions[1], Pos{x: 2, y: 994})
    assert.Equal(t, positions[2], Pos{x: 2, y: 990})
    assert.Equal(t, positions[3], Pos{x: 6, y: 998})
    assert.Equal(t, positions[4], Pos{x: 6, y: 994})
    assert.Equal(t, positions[5], Pos{x: 6, y: 990})
    assert.Equal(t, positions[6], Pos{x: 10, y: 998})
    assert.Equal(t, positions[7], Pos{x: 10, y: 994})
    assert.Equal(t, positions[8], Pos{x: 14, y: 998})
    assert.Equal(t, positions[9], Pos{x: 14, y: 994})
    assert.Equal(t, positions[10], Pos{x: 14, y: 990})
    assert.Equal(t, positions[11], Pos{x: 18, y: 998})
    assert.Equal(t, positions[12], Pos{x: 18, y: 994})
}
