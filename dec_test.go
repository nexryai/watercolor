package watercolor

import (
	"image"
	"image/png"
	"os"
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
