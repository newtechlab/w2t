// Command w2t changes white pixels in images to be transparent
package main

import (
	"flag"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"sync"
)

var (
	fInPlace bool
)

func init() {
	flag.BoolVar(&fInPlace, "replace", false, "replace the images")
}

func main() {
	flag.Parse()
	queue := make(chan struct{}, 10)
	wg := &sync.WaitGroup{}
	for _, p := range flag.Args() {
		wg.Add(1)
		go do(queue, wg, p)
	}
	wg.Wait()
}

func do(queue chan struct{}, wg *sync.WaitGroup, p string) {
	queue <- struct{}{}
	defer wg.Done()
	defer func() {
		<-queue
	}()

	f, err := os.Open(p)
	check(err)
	defer f.Close()

	img, _, err := image.Decode(f)
	check(err)
	f.Close()
	bnd := img.Bounds()

	// for now always save as colored and transparent png
	out := image.NewNRGBA(bnd)

	for x := 0; x < bnd.Max.X; x++ {
		for y := 0; y < bnd.Max.Y; y++ {
			col := img.At(x, y)
			r, g, b, a := col.RGBA()
			if true || a == 0xffff && r == g && r == b {
				rr := uint16(r)
				out.Set(x, y, color.NRGBA{0, 0, 0, 255 - uint8(rr>>8)})
			} else {
				out.Set(x, y, col)
			}
		}
	}

	if fInPlace {
		f, err = os.Create(p)
		check(err)
		check(png.Encode(f, out))
	} else {
		f, err = os.Create(p + ".mod")
		check(err)
		check(png.Encode(f, out))
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
