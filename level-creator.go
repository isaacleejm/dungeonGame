package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var Layouts = struct {
	Layout1 []string
	Layout2 []string
	Layout3 []string
	Layout4 []string
	Layout5 []string
	Layout6 []string
	Layout7 []string
	Layout8 []string
	Layout9 []string
	Layout10 []string
} {
	Layout1: []string{
		"####################",
		"#                  #",
		"#   ###      ###   #",
		"#   ###      ###   #",
		"#                  #",
		"#       ####       #",
		"#                  #",
		"#  ###        ###  #",
		"#  ###        ###  #",
		"#                  #",
		"#       ####       #",
		"#                  #",
		"#   ###      ###   #",
		"#                  #",
		"####################",
	},

	Layout2: []string{
		"####################",
		"#                  #",
		"#                  #",
		"#      ###         #",
		"#      ###         #",
		"#                  #",
		"#   ###            #",
		"#   ###            #",
		"#            ###   #",
		"#            ###   #",
		"#                  #",
		"#         ###      #",
		"#         ###      #",
		"#                  #",
		"####################",
	},

	Layout3: []string{
		"####################",
		"#                  #",
		"#  ####      ####  #",
		"#  ####      ####  #",
		"#                  #",
		"######        ######",
		"#                  #",
		"#        ####      #",
		"#        ####      #",
		"#                  #",
		"######        ######",
		"#                  #",
		"#  ####      ####  #",
		"#                  #",
		"####################",
	},

	Layout4: []string{
		"####################",
		"#                  #",
		"#  ####            #",
		"#  ####            #",
		"#                  #",
		"#            ####  #",
		"#            ####  #",
		"#                  #",
		"#      ####        #",
		"#      ####        #",
		"#                  #",
		"#        ####      #",
		"#        ####      #",
		"#                  #",
		"####################",
	},

	Layout5: []string{
		"####################",
		"#                  #",
		"#   ##        ##   #",
		"#   ##        ##   #",
		"#                  #",
		"#   ##        ##   #",
		"#                  #",
		"#        ####      #",
		"#        ####      #",
		"#                  #",
		"#   ##        ##   #",
		"#                  #",
		"#   ##        ##   #",
		"#                  #",
		"####################",
	},

	Layout6: []string{
		"####################",
		"#                  #",
		"#  ####            #",
		"#  ####            #",
		"#  ####            #",
		"#                  #",
		"#          ####    #",
		"#          ####    #",
		"#          ####    #",
		"#                  #",
		"#    ####          #",
		"#    ####          #",
		"#    ####          #",
		"#                  #",
		"####################",
	},

	Layout7: []string{
		"####################",
		"#                  #",
		"#   ####           #",
		"#   ####           #",
		"#                  #",
		"#        ####      #",
		"#        ####      #",
		"#                  #",
		"#            ####  #",
		"#            ####  #",
		"#                  #",
		"#        ####      #",
		"#        ####      #",
		"#                  #",
		"####################",
	},

	Layout8: []string{
		"####################",
		"#                  #",
		"#  ###       ###   #",
		"#  ###       ###   #",
		"#                  #",
		"#       ###        #",
		"#       ###        #",
		"#                  #",
		"#       ###        #",
		"#       ###        #",
		"#                  #",
		"#  ###       ###   #",
		"#  ###       ###   #",
		"#                  #",
		"####################",
	},

	Layout9: []string{
		"####################",
		"#                  #",
		"#  ####            #",
		"#  ####            #",
		"#       ####       #",
		"#       ####       #",
		"#                  #",
		"#            ####  #",
		"#            ####  #",
		"#                  #",
		"#  ####            #",
		"#  ####            #",
		"#       ####       #",
		"#                  #",
		"####################",
	},

	Layout10: []string{
		"####################",
		"#                  #",
		"#   ###    ###     #",
		"#   ###    ###     #",
		"#                  #",
		"#      ####        #",
		"#      ####        #",
		"#                  #",
		"#        ####      #",
		"#        ####      #",
		"#                  #",
		"#     ###    ###   #",
		"#     ###    ###   #",
		"#                  #",
		"####################",
	},
}

var Backgrounds = struct {
	Blue   color.RGBA
	Red    color.RGBA
	Green  color.RGBA
	Purple color.RGBA
	Gray   color.RGBA
}{
	Blue: color.RGBA{
		R: 30, G: 30, B: 46, A: 255,
	},
	Red: color.RGBA{
		R: 50, G: 25, B: 25, A: 255,
	},
	Green: color.RGBA{
		R: 25, G: 50, B: 30, A: 255,
	},
	Purple: color.RGBA{
		R: 45, G: 25, B: 55, A: 255,
	},
	Gray: color.RGBA{
		R: 45, G: 45, B: 45, A: 255,
	},
}

type LevelMap struct {
	Blocks     []*Block
	Background color.RGBA
}

func BlocksFromLayout(layout []string, Sprite *ebiten.Image, background color.RGBA) LevelMap {
	blocks := []*Block{}

	const blockSize = 50

	for y, row := range layout {
		for x, tile := range row {
			if tile == '#' {
				blocks = append(blocks, &Block{
					Pos: Vector2{
						X: float64(x * blockSize),
						Y: float64(y * blockSize),
					},
					Sprite: Sprite,
				})
			}
		}
	}

	return LevelMap{
		Blocks:     blocks,
		Background: background,
	}
}
