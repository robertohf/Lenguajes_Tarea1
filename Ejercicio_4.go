package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"os"
)

func main() {
	bitmap, err := openImage("lena.bmp")
	if err != nil {
		fmt.Println(err)
	}

	bounds := bitmap.Bounds()
	//w, h := bounds.Max.X, bounds.Max.Y

	imgSet := image.NewRGBA(image.Rect(0, 0, 128, 128))

	w, h := bounds.Max.X/128, bounds.Max.X/128

	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			pixel := bitmap.At(x*w, y*h)
			imgSet.Set(x, y, pixel)
		}
	}

	outfile, err := os.Create("lena_Redux.bmp")
	if err != nil {
		fmt.Println(err)
	}

	defer outfile.Close()

	bmp.Encode(outfile, imgSet)

}

func openImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return bmp.Decode(f)
}
