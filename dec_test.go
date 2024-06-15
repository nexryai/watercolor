package watercolor

import (
	"fmt"
	"github.com/nexryai/apng"
	"image"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

func savePNG(filename string, img *image.RGBA) error {
	// ファイルを作成またはオープン
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// PNGエンコードしてファイルに書き込む
	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

func blendFrames(prevFrame *image.RGBA, currFrame *image.Image, offsetX int, offsetY int) *image.RGBA {
	// Create a new image with the same size as the previous frame
	blended := image.NewRGBA(prevFrame.Bounds())

	// Copy the previous frame to the new image if the dispose operation is not OpBackground
	draw.Draw(blended, blended.Bounds(), prevFrame, image.Point{}, draw.Src)

	// Determine the blend operation
	op := draw.Over

	// Draw the current frame onto the new image
	draw.Draw(blended, blended.Bounds(), *currFrame, image.Point{-offsetX, -offsetY}, op)

	return blended
}

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

func TestDecApng(t *testing.T) {
	// Open our animated PNG file
	f, err := os.Open("testdata/apng1.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var baseFrame *image.RGBA

	_, err = apng.DecodeAll(f, func(f *apng.FrameHookArgs) error {
		fmt.Printf("Frame %d | Delay:%v\n", f.Num, f.Delay)
		img := f.Buffer

		if f.Num == 0 {
			baseFrame = convertToRGBA(img)
			return nil
		}

		var rgba *image.RGBA
		if f.BlendOp == apng.BLEND_OP_OVER {
			rgba = blendFrames(baseFrame, img, f.OffsetX, f.OffsetY)
		} else {
			rgba = convertToRGBA(img)
		}

		fmt.Printf("DisposeOp: %d\n", f.DisposeOp)
		switch f.DisposeOp {
		case apng.DISPOSE_OP_NONE:
			baseFrame = rgba
		case apng.DISPOSE_OP_BACKGROUND:
			// Clear the frame
			draw.Draw(rgba, rgba.Bounds(), image.Transparent, image.Point{}, draw.Src)
		case apng.DISPOSE_OP_PREVIOUS:
			// 一つ前のフレームをそのまま使用するため何もしない
		}

		// encode to png
		err := savePNG(fmt.Sprintf("testout/apng1_%d.png", f.Num), rgba)
		if err != nil {
			t.Fatal(err)
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}
