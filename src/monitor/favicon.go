package main

import (
	"image"
	"image/color"
	"image/draw"
)

func Favicon() image.Image {
	icon := image.NewRGBA(image.Rect(0, 0, 16, 16))
	blue := color.RGBA{0, 0, 255, 255}
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(icon, image.Rect(6, 7, 10, 15), &image.Uniform{black}, image.ZP, draw.Src)
	draw.Draw(icon, image.Rect(6, 14, 14, 15), &image.Uniform{black}, image.ZP, draw.Src)
	draw.Draw(icon, image.Rect(4, 7, 10, 8), &image.Uniform{black}, image.ZP, draw.Src)
	draw.Draw(icon, image.Rect(7, 2, 10, 5), &image.Uniform{blue}, image.ZP, draw.Src)
	return icon
}
