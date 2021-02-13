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
	levelWidth   = 64
	levelHeight  = 16
)

var (
	cameraPosX = 0.0
	cameraPosY = 0.0

	playerPosX = 2.0
	playerPosY = 2.0
	playerVX   = 0.0
	playerVY   = 0.0
)

var levelString = "################################################################" +
	"#..............................................................#" +
	"#...............###.........#########................######....#" +
	"#........###..............................#####................#" +
	"#.....................###............#####.....................#" +
	"#..............................................######..........#" +
	"#..............................................................#" +
	"#..............................................................#" +
	"###########################.###############.....################" +
	"#.........................#.##................###..............#" +
	"#.........................#.#................###...............#" +
	"#..................########.#.............###..................#" +
	"#..................#........#..........###.....................#" +
	"#..................#.########.......###........................#" +
	"#..................#.............###...........................#" +
	"################################################################"

func getCell(x, y int) byte {
	if x < 0 || x >= levelWidth || y < 0 || y >= levelHeight {
		return ' '
	}
	i := y*levelWidth + x
	return levelString[i]
}

type Game struct{}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0})
	tileWidth := 64
	tileHeight := 64
	visibleTileX := (screenWidth / tileWidth)
	visibleTileY := (screenHeight / tileHeight)

	cameraPosX = playerPosX
	cameraPosY = playerPosY

	offsetX := clampFloat(cameraPosX-float64(visibleTileX)/2.0, 0, float64(levelWidth-visibleTileX))
	offsetY := clampFloat(cameraPosY-float64(visibleTileY)/2.0, 0, float64(levelHeight-visibleTileY))

	tileOffsetX := (offsetX - float64(int(offsetX))) * float64(tileWidth)
	tileOffsetY := (offsetY - float64(int(offsetY))) * float64(tileHeight)

	for x := -1; x < visibleTileX+1; x++ {
		for y := -1; y < visibleTileY+2; y++ {
			c := getCell(x+int(offsetX), y+int(offsetY))
			// fmt.Printf("%c\n", c)
			if c == '.' {
				ebitenutil.DrawRect(screen, float64(x*tileWidth)-tileOffsetX, float64(y*tileHeight)-tileOffsetY, float64(tileWidth), float64(tileHeight), color.RGBA{0, 0xCD, 0xDF, 0xFF})
			}
			if c == '#' {
				ebitenutil.DrawRect(screen, float64(x*tileWidth)-tileOffsetX, float64(y*tileHeight)-tileOffsetY, float64(tileWidth), float64(tileHeight), color.RGBA{0x6D, 0x3B, 0x11, 0xFF})
			}
		}
	}
	//draw player
	// ebitenutil.DrawRect(screen, (playerPosX-offsetX)*float64(tileWidth)-tileOffsetX, (playerPosY-offsetY)*float64(tileHeight)-tileOffsetY, float64(tileWidth), float64(tileHeight), color.RGBA{0x15, 0x6D, 0x11, 0xFF})
	ebitenutil.DrawRect(screen, (playerPosX-offsetX)*float64(tileWidth), (playerPosY-offsetY)*float64(tileHeight), float64(tileWidth), float64(tileHeight), color.RGBA{0x15, 0x6D, 0x11, 0xFF})

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

var onTheGround = false

func (g *Game) Update() error {
	elapsed := 1.0 / 60.0
	// fmt.Printf("On the ground : %v\n", onTheGround)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		//player on ground difference there
		if onTheGround {
			playerVX += (-25.0) * elapsed
		} else {
			playerVX += (-15.0) * elapsed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if onTheGround {
			playerVX += (25.0) * elapsed
		} else {
			playerVX += (15.0) * elapsed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if playerVY == 0 && onTheGround {
			playerVY -= 12.0
		}
	}
	playerVY += 20.0 * elapsed

	//Drag if on the ground
	if onTheGround {
		playerVX += -3.0 * playerVX * elapsed
		if math.Abs(playerVX) < 0.01 {
			playerVX = 0.0
		}
	}

	playerVX = clampFloat(playerVX, -10.0, 10.0)
	playerVY = clampFloat(playerVY, -100, 100)
	newPlayerPosX := playerPosX + playerVX*elapsed
	newPlayerPosY := playerPosY + playerVY*elapsed

	onTheGround = false
	//check for collision
	if playerVX <= 0 {
		if getCell(int(newPlayerPosX), int(playerPosY)) != '.' || getCell(int(newPlayerPosX), int(playerPosY+0.99)) != '.' {
			newPlayerPosX = float64(int(newPlayerPosX) + 1)
			playerVX = 0
		}
	}
	if playerVX >= 0 {
		if getCell(int(newPlayerPosX+0.99), int(playerPosY)) != '.' || getCell(int(newPlayerPosX+0.99), int(playerPosY+0.99)) != '.' {
			newPlayerPosX = float64(int(newPlayerPosX))
			playerVX = 0
		}
	}
	if playerVY <= 0 {
		if getCell(int(newPlayerPosX), int(newPlayerPosY)) != '.' || getCell(int(newPlayerPosX+0.99), int(newPlayerPosY)) != '.' {
			newPlayerPosY = float64(int(newPlayerPosY) + 1)
			playerVY = 0
		}
	}
	if playerVY >= 0 {
		if getCell(int(newPlayerPosX), int(newPlayerPosY+0.99)) != '.' || getCell(int(newPlayerPosX+0.99), int(newPlayerPosY+0.99)) != '.' {
			newPlayerPosY = float64(int(newPlayerPosY))
			playerVY = 0
			onTheGround = true
		}
	}

	playerPosX = newPlayerPosX
	playerPosY = newPlayerPosY
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (sW, sH int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
