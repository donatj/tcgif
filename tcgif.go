package tcgif

import (
	"image"
	"image/color"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"sort"
)

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

type GIFMaker struct {
	frameDelay int
	finalDelay int

	frameLimit uint

	backfill bool
	popsort  bool
}

type Option func(*GIFMaker)

func WithFrameDelay(delay int) Option {
	return func(g *GIFMaker) {
		g.frameDelay = delay
	}
}

func WithFinalDelay(delay int) Option {
	return func(g *GIFMaker) {
		g.finalDelay = delay
	}
}

func WithFrameLimit(limit uint) Option {
	return func(g *GIFMaker) {
		g.frameLimit = limit
	}
}

func WithBackfill(backfill bool) Option {
	return func(g *GIFMaker) {
		g.backfill = backfill
	}
}

func WithPopularitySort(popsort bool) Option {
	return func(g *GIFMaker) {
		g.popsort = popsort
	}
}

func NewGIFMaker(opts ...Option) *GIFMaker {
	gm := &GIFMaker{
		frameDelay: 2,
		finalDelay: 300,
		frameLimit: 0,
		backfill:   true,
		popsort:    true,
	}

	for _, opt := range opts {
		opt(gm)
	}

	return gm
}

func (gm *GIFMaker) MakeGIF(img image.Image) (*gif.GIF, error) {
	b := img.Bounds()
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

	if gm.popsort {
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
	if gm.frameLimit != 0 && int(gm.frameLimit) < limitSeglen {
		limitSeglen = int(gm.frameLimit)
	}

	g := &gif.GIF{}
	for i := 0; i < limitSeglen; i++ {
		pimg := image.NewPaletted(b, color.Palette{})
		// Add trasparency first so it's used as the matte color
		pimg.Palette = append(pimg.Palette, color.Transparent)
		g.Image = append(g.Image, pimg)

		for _, ch := range segments[i] {
			pimg.Palette = append(pimg.Palette, ch.C)
			ind := pimg.Palette.Index(ch.C)

			for _, ccoord := range ch.Coords {
				pimg.SetColorIndex(ccoord.X, ccoord.Y, uint8(ind))
			}
		}

		if gm.backfill {
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
		g.Delay[i] = gm.frameDelay
	}

	g.Delay[len(g.Delay)-1] = gm.finalDelay

	return g, nil
}
