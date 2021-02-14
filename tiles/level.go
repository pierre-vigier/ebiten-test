package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Level struct {
	lWidth  int
	lHeight int
	lMap    string
}

func NewLevel() *Level {
	return &Level{
		lMap: "################################################################" +
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
			"################################################################",
		lHeight: 16,
		lWidth:  64,
	}
}

func (l *Level) getCell(x, y int) byte {
	if x < 0 || x >= l.lWidth || y < 0 || y >= l.lHeight {
		return ' '
	}
	i := y*l.lWidth + x
	return l.lMap[i]
}

func (l *Level) Width() int {
	return l.lWidth
}

func (l *Level) Height() int {
	return l.lHeight
}

func (l *Level) Draw(screen *ebiten.Image, offsetX, offsetY float64) {
	tileOffsetX := (offsetX - float64(int(offsetX))) * float64(tileWidth)
	tileOffsetY := (offsetY - float64(int(offsetY))) * float64(tileHeight)

	for x := -1; x < visibleTileX+1; x++ {
		for y := -1; y < visibleTileY+2; y++ {
			c := l.getCell(x+int(offsetX), y+int(offsetY))
			// fmt.Printf("%c\n", c)
			if c == '.' {
				ebitenutil.DrawRect(screen, float64(x*tileWidth)-tileOffsetX, float64(y*tileHeight)-tileOffsetY, float64(tileWidth), float64(tileHeight), color.RGBA{0, 0xCD, 0xDF, 0xFF})
			}
			if c == '#' {
				ebitenutil.DrawRect(screen, float64(x*tileWidth)-tileOffsetX, float64(y*tileHeight)-tileOffsetY, float64(tileWidth), float64(tileHeight), color.RGBA{0x6D, 0x3B, 0x11, 0xFF})
			}
		}
	}
}
