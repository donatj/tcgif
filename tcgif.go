package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"sort"
)

var (
	frameLimit = flag.Uint("framelimit", 0, "max number of frames. 0 = unlimited")
	backfill   = flag.Bool("backfill", true, "backfill still missing pixels with closest color")
	popsort    = flag.Bool("sort", true, "sort colors by popularity")
)

func init() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		fmt.Println("requires one image as input")
		os.Exit(1)
	}
}

type coord struct {
	X, Y int
}

type colorCount struct {
	C      color.Color
	Coords []coord
}

type colorCountList []colorCount

func (p colorCountList) Len() int           { return len(p) }
func (p colorCountList) Less(i, j int) bool { return len(p[i].Coords) < len(p[j].Coords) }
func (p colorCountList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	b := img.Bounds()
	g := &gif.GIF{}

	colormap := make(map[color.Color][]coord)

	for y := 0; y <= b.Max.Y; y++ {
		for x := 0; x <= b.Max.X; x++ {
			c := img.At(x, y)
			colormap[c] = append(colormap[c], coord{x, y})
		}
	}

	colorhisto := make(colorCountList, 0)
	for c, e := range colormap {
		colorhisto = append(colorhisto, colorCount{c, e})
	}

	if *popsort {
		sort.Sort(sort.Reverse(colorhisto))
	}

	seglen := (len(colorhisto) / 254) + 1
	segments := make([]colorCountList, seglen)

	x := 0
	for _, xxx := range colorhisto {
		n := x / 254 //integer division
		segments[n] = append(segments[n], xxx)

		x++
	}

	limitSeglen := seglen
	if *frameLimit != 0 && int(*frameLimit) < limitSeglen {
		limitSeglen = int(*frameLimit)
	}

	for i := 0; i < limitSeglen; i++ {
		pimg := image.NewPaletted(b, color.Palette{})
		pimg.Palette = append(pimg.Palette, color.Transparent)
		g.Image = append(g.Image, pimg)

		for _, ch := range segments[i] {
			pimg.Palette = append(pimg.Palette, ch.C)
			ind := pimg.Palette.Index(ch.C)

			for _, ccoord := range ch.Coords {
				pimg.SetColorIndex(ccoord.X, ccoord.Y, uint8(ind))
			}
		}

		if *backfill {
			for j := i + 1; j < seglen; j++ {
				for _, ch := range segments[j] {
					ind := pimg.Palette.Index(ch.C)

					for _, ccoord := range ch.Coords {
						pimg.SetColorIndex(ccoord.X, ccoord.Y, uint8(ind))
					}
				}
			}
		}
	}

	g.Delay = make([]int, len(g.Image))
	for i := range g.Delay {
		g.Delay[i] = 0
	}

	out, err := os.Create("out.gif")
	if err != nil {
		log.Fatal(err)
	}

	err = gif.EncodeAll(out, g)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Output to: out.gif")
	fmt.Printf("Conatins %d frames.\n", len(g.Image))
}
