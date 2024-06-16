package watercolor

/*
   #cgo LDFLAGS: -lwebp -lwebpmux
   #cgo darwin pkg-config: libwebp
   #include <webp/encode.h>
   #include <webp/mux.h>
*/
import "C"

import (
	"bytes"
	"fmt"
	"github.com/nexryai/apng"
	"golang.org/x/image/draw"
	"image"
	"unsafe"
)

// BLEND_OP_OVERの場合に前のフレームとのブレンドする
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

// Goのimage.Imageをlibwebp.WebPPictureに変換
func imageToWebPPicture(img *image.RGBA, scale float32, width int, height int, xOffset int, yOffset int) C.WebPPicture {
	bounds := (*img).Bounds()
	fmt.Printf("Dx: %v Dy: %v\n", bounds.Dx(), bounds.Dy())

	var pic C.WebPPicture
	C.WebPPictureInit(&pic)

	pic.width = C.int(width)
	pic.height = C.int(height)
	pic.use_argb = 1

	// リサイズ後の画像を格納するためのRGBAイメージを作成
	rgbaImg := image.NewRGBA(image.Rect(0, 0, width, height))

	if scale != 1 {
		// リサイズ
		newWidth := int(float32((*img).Bounds().Dx()) * scale)
		newHeight := int(float32((*img).Bounds().Dy()) * scale)

		//xOffset = int(float32(xOffset) * scale)
		//yOffset = int(float32(yOffset) * scale)

		fmt.Printf("newWidth: %v newHeight: %v\n", newWidth, newHeight)
		draw.ApproxBiLinear.Scale(rgbaImg, image.Rect(xOffset, yOffset, newWidth, newHeight), img, bounds, draw.Src, nil)
	} else {
		// リサイズしない場合、そのままコピー
		rgbaImg = img
	}

	// WebPにエンコード
	C.WebPPictureImportRGBA(&pic, (*C.uint8_t)(unsafe.Pointer(&rgbaImg.Pix[0])), C.int(rgbaImg.Stride))

	return pic
}

func apngToWebP(imgPtr *[]byte, width int, height int) (*[]byte, error) {
	// libwebpの初期化
	buffer := bytes.NewBuffer(*imgPtr)

	// Skip the first 8 bytes (PNG signature)
	buffer.Next(8)

	// Skip chunk type (8 bytes)
	buffer.Next(8)

	originalWidth := readInt32(buffer)
	originalHeight := readInt32(buffer)

	// 0 divide check
	if originalWidth == 0 || originalHeight == 0 {
		return nil, fmt.Errorf("invalid image size")
	}

	// 0の場合はリサイズしないで元のサイズを使用
	if width == 0 {
		width = originalWidth
	}
	if height == 0 {
		height = originalHeight
	}

	fmt.Printf("originalWidth: %d, originalHeight: %d\n", originalWidth, originalHeight)

	var animConfig C.WebPAnimEncoderOptions
	C.WebPAnimEncoderOptionsInit(&animConfig)
	animEncoder := C.WebPAnimEncoderNew(C.int(width), C.int(height), &animConfig)

	scale := float32(height) / float32(originalHeight)
	fmt.Printf("scale: %v\n", scale)

	// BLEND_OP_OVERの場合、前のフレームとのブレンドが必要
	var baseFrame *image.RGBA

	_, err := apng.DecodeAll(bytes.NewReader(*imgPtr), func(f *apng.FrameHookArgs) error {
		// 以下フレームごとに実行される処理
		if f.Num == 0 {
			// 最初のフレームをRGBAに変換してbaseFrameに保存
			baseFrame = convertToRGBA(f.Buffer)

			// 最初のフレームはスキップ
			return nil
		}

		var rgba *image.RGBA
		if f.BlendOp == apng.BLEND_OP_OVER {
			rgba = blendFrames(baseFrame, f.Buffer, f.OffsetX, f.OffsetY)
		} else {
			rgba = convertToRGBA(f.Buffer)
		}

		// DisposeOpによって処理を変更
		switch f.DisposeOp {
		case apng.DISPOSE_OP_NONE:
			baseFrame = rgba
		case apng.DISPOSE_OP_BACKGROUND:
			// Clear the frame
			draw.Draw(rgba, rgba.Bounds(), image.Transparent, image.Point{}, draw.Src)
		case apng.DISPOSE_OP_PREVIOUS:
			// 一つ前のフレームをそのまま使用するため何もしない
		}

		// webpとしてエンコード
		pic := imageToWebPPicture(rgba, scale, width, height, f.OffsetX, f.OffsetY)

		// リサイズ
		C.WebPPictureRescale(&pic, C.int(width), C.int(height))

		timeStamp := int(float32(f.Num) * f.Delay * 1000)
		fmt.Printf("timeStamp: %d\n", timeStamp)

		// Animated WebPのフレームとして追加
		result := C.int(C.WebPAnimEncoderAdd(animEncoder, &pic, C.int(timeStamp), nil))
		if result == 0 {
			// animEncoderの解放
			C.WebPPictureFree(&pic)
			return fmt.Errorf("WebPAnimEncoderAdd failed")
		}

		// Cleanup
		C.WebPPictureFree(&pic)

		return nil
	})

	if err != nil {
		C.WebPAnimEncoderDelete(animEncoder)
		return nil, err
	}

	// 書き込み
	var webpData C.WebPData
	C.WebPDataInit(&webpData)
	C.WebPAnimEncoderAssemble(animEncoder, &webpData)
	webpBytes := C.GoBytes(unsafe.Pointer(webpData.bytes), C.int(webpData.size))

	// animEncoderの解放
	C.WebPDataClear(&webpData)
	C.WebPAnimEncoderDelete(animEncoder)

	return &webpBytes, nil
}
