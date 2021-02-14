package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	// "github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 1280
	screenHeight = 960
	tileWidth    = 64
	tileHeight   = 64
)

var (
	cameraPosX   = 0.0
	cameraPosY   = 0.0
	visibleTileX = (screenWidth / tileWidth)
	visibleTileY = (screenHeight / tileHeight)
)

type Vector struct {
	X, Y float64
}

type Player struct {
	Pos         Vector
	Vel         Vector
	OnTheGround bool
}

func NewPlayer() *Player {
	return &Player{
		Pos: Vector{
			X: 2.0,
			Y: 2.0,
		},
		Vel: Vector{
			X: 0.0,
			Y: 0.0,
		},
		OnTheGround: false,
	}
}

func (p *Player) Draw(screen *ebiten.Image, offsetX, offsetY float64) {
	ebitenutil.DrawRect(screen, (p.Pos.X-offsetX)*float64(tileWidth), (p.Pos.Y-offsetY)*float64(tileHeight), float64(tileWidth), float64(tileHeight), color.RGBA{0x15, 0x6D, 0x11, 0xFF})
}

func (p *Player) MoveLeft(elapsed float64) {
	//player on ground difference there
	if p.OnTheGround {
		p.Vel.X += (-25.0) * elapsed
	} else {
		p.Vel.X += (-15.0) * elapsed
	}
}

func (p *Player) MoveRight(elapsed float64) {
	//player on ground difference there
	if p.OnTheGround {
		p.Vel.X += (25.0) * elapsed
	} else {
		p.Vel.X += (15.0) * elapsed
	}
}

func (p *Player) Jump(elapsed float64) {
	if p.Vel.Y == 0 && p.OnTheGround {
		p.Vel.Y -= 12.0
	}
}

func (p *Player) Apply(gravity, elapsed float64) {
	p.Vel.Y += gravity * elapsed
	if p.OnTheGround {
		p.Vel.X += -3.0 * p.Vel.X * elapsed
		if math.Abs(p.Vel.X) < 0.01 {
			p.Vel.X = 0.0
		}
	}
	p.Vel.X = clampFloat(p.Vel.X, -10.0, 10.0)
	p.Vel.Y = clampFloat(p.Vel.Y, -100, 100)
}

type Game struct {
	level  *Level
	player *Player
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0})

	cameraPosX = g.player.Pos.X
	cameraPosY = g.player.Pos.Y

	offsetX := clampFloat(cameraPosX-float64(visibleTileX)/2.0, 0, float64(g.level.Width()-visibleTileX))
	offsetY := clampFloat(cameraPosY-float64(visibleTileY)/2.0, 0, float64(g.level.Height()-visibleTileY))

	g.level.Draw(screen, offsetX, offsetY)
	//draw player
	g.player.Draw(screen, offsetX, offsetY)

}

func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func (g *Game) Update() error {
	elapsed := 1.0 / 60.0
	// fmt.Printf("On the ground : %v\n", onTheGround)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.MoveLeft(elapsed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.MoveRight(elapsed)
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.player.Jump(elapsed)
	}

	g.player.Apply(20.0, elapsed)

	//check collision
	newPlayerPosX := g.player.Pos.X + g.player.Vel.X*elapsed
	newPlayerPosY := g.player.Pos.Y + g.player.Vel.Y*elapsed

	g.player.OnTheGround = false
	//check for collision
	if g.player.Vel.X <= 0 {
		if g.level.getCell(int(newPlayerPosX), int(g.player.Pos.Y)) != '.' || g.level.getCell(int(newPlayerPosX), int(g.player.Pos.Y+0.99)) != '.' {
			newPlayerPosX = float64(int(newPlayerPosX) + 1)
			g.player.Vel.X = 0
		}
	}
	if g.player.Vel.X >= 0 {
		if g.level.getCell(int(newPlayerPosX+0.99), int(g.player.Pos.Y)) != '.' || g.level.getCell(int(newPlayerPosX+0.99), int(g.player.Pos.Y+0.99)) != '.' {
			newPlayerPosX = float64(int(newPlayerPosX))
			g.player.Vel.X = 0
		}
	}
	if g.player.Vel.Y <= 0 {
		if g.level.getCell(int(newPlayerPosX), int(newPlayerPosY)) != '.' || g.level.getCell(int(newPlayerPosX+0.99), int(newPlayerPosY)) != '.' {
			newPlayerPosY = float64(int(newPlayerPosY) + 1)
			g.player.Vel.Y = 0
		}
	}
	if g.player.Vel.Y >= 0 {
		if g.level.getCell(int(newPlayerPosX), int(newPlayerPosY+0.99)) != '.' || g.level.getCell(int(newPlayerPosX+0.99), int(newPlayerPosY+0.99)) != '.' {
			newPlayerPosY = float64(int(newPlayerPosY))
			g.player.Vel.Y = 0
			g.player.OnTheGround = true
		}
	}

	g.player.Pos.X = newPlayerPosX
	g.player.Pos.Y = newPlayerPosY
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (sW, sH int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	game := &Game{
		level:  NewLevel(),
		player: NewPlayer(),
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
