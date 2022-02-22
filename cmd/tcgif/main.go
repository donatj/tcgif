package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"

	"image/gif"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/donatj/tcgif"
)

var (
	// See: https://www.biphelps.com/blog/The-Fastest-GIF-Does-Not-Exist
	frameDelay = flag.Int("delay", 2, "frame delay in multiples of 10ms. 2 is fastest for historical reasons")
	finalDelay = flag.Int("final-delay", 300, "frame delay in on final frame")

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

func main() {
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	gm := tcgif.NewGIFMaker(
		tcgif.WithFrameDelay(*frameDelay),
		tcgif.WithFinalDelay(*finalDelay),
		tcgif.WithFrameLimit(*frameLimit),
		tcgif.WithBackfill(*backfill),
		tcgif.WithPopularitySort(*popsort),
	)

	g, err := gm.MakeGIF(img)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create("out.gif")
	if err != nil {
		log.Fatal(err)
	}

	err = gif.EncodeAll(out, g)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Output %d frames to: out.gif\n", len(g.Image))
}
