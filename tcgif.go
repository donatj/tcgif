package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"os"
)

func init() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("requires one jpeg as input")
		os.Exit(1)
	}
}

func main() {
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	b := img.Bounds()

	g := &gif.GIF{}

	j := 0
	var pimg *image.Paletted

	for x := 0; x <= b.Max.X; x++ {
		for y := 0; y <= b.Max.Y; y++ {

			ind := j % 254

			if ind == 0 {
				pimg = image.NewPaletted(b, color.Palette{})
				pimg.Palette = append(pimg.Palette, color.Transparent)
				g.Image = append(g.Image, pimg)
			}

			c := img.At(x, y)
			//			fmt.Println(ind, x, y)
			//			pimg.Palette[ind] = c

			ci := -1
			for cc, pp := range pimg.Palette {
				if colorEq(c, pp) {
					ci = cc
				}
			}

			if ci == -1 {
				pimg.Palette = append(pimg.Palette, c)
				ci = len(pimg.Palette) - 1
			} else {
				j -= 1
				fmt.Println(".")
			}

			pimg.SetColorIndex(x, y, uint8(ci))

			j++
		}
	}

	g.Delay = make([]int, len(g.Image))
	for i, _ := range g.Delay {
		g.Delay[i] = 0
	}

	fmt.Println(j)

	out, err := os.Create("out.gif")
	if err != nil {
		panic(err)
	}

	err = gif.EncodeAll(out, g)
	if err != nil {
		panic(err)
	}

	fmt.Println("Output to: out.gif")
	fmt.Printf("Conatins %d frames.\n", len(g.Image))
}

func colorEq(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()

	if ar == br && ag == bg && ab == bb && aa == ba {
		return true
	}

	return false
}
