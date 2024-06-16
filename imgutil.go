package watercolor

import "image"

func convertToRGBA(img *image.Image) *image.RGBA {
	bounds := (*img).Bounds()
	rgba := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, (*img).At(x, y))
		}
	}

	return rgba
}
