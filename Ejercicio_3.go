package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"os"
)

func main() {
	bitmap, err := openImage("lena.bmp")
	if err != nil {
		fmt.Println(err)
	}

	bounds := bitmap.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	imgSet := image.NewRGBA(bounds)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			oldPixel := bitmap.At(x, y)
			r, g, b, _ := oldPixel.RGBA()
			//CAMBIA ESTA MIERDA AQUI ABAJO
			avg := (r + g + b) / 3 //ESTOOOO!!!
			pixel := color.Gray{uint8(avg / 256)}
			imgSet.Set(x, y, pixel)
		}
	}

	outfile, err := os.Create("lena_GrayScale.bmp")
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
