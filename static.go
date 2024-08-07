package watercolor

import (
	"bytes"
	"fmt"
	avif "github.com/nexryai/goavif"
	ico "github.com/ur65/go-ico"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
	"image"
	"image/jpeg"
	"image/png"
)

var (
	ErrInvalidParam            = fmt.Errorf("invalid parameter")
	ErrImageIsNil              = fmt.Errorf("image is nil")
	ErrImageFormatUnknown      = fmt.Errorf("image format is unknown")
	ErrImageFormatNotSupported = fmt.Errorf("image format is not supported")
	ErrFailedToDecodeJPEG      = fmt.Errorf("failed to decode as jpeg")
	ErrFailedToDecodePNG       = fmt.Errorf("failed to decode as png")
	ErrFailedToDecodeWebP      = fmt.Errorf("failed to decode as webp")
	ErrFailedToDecodeAVIF      = fmt.Errorf("failed to decode as avif")
	ErrFailedToEncodeWebP      = fmt.Errorf("failed to encode as webp")
	ErrFailedToEncodeAVIF      = fmt.Errorf("failed to encode as avif")
)

func ProcessStaticImage(data *[]byte, targetImage *TargetImage) (*[]byte, error) {
	if targetImage == nil {
		return nil, ErrInvalidParam
	} else if data == nil {
		return nil, ErrImageIsNil
	}

	format := detectImageFormat(data)
	if format == ImageFormatUnknown {
		return nil, ErrImageFormatUnknown
	} else if format == ImageFormatAnimatedPNG {
		// Staticな画像として処理
		format = ImageFormatPNG
	}

	imageReader := bytes.NewReader(*data)
	var decodedImage image.Image
	var err error

	switch format {
	case ImageFormatJPEG:
		decodedImage, err = jpeg.Decode(imageReader)
		if err != nil {
			return nil, ErrFailedToDecodeJPEG
		}
	case ImageFormatPNG:
		decodedImage, err = png.Decode(imageReader)
		if err != nil {
			return nil, ErrFailedToDecodePNG
		}
	case ImageFormatWebP:
		decodedImage, err = webp.Decode(imageReader)
		if err != nil {
			return nil, ErrFailedToDecodeWebP
		}
	case ImageFormatAVIF:
		decodedImage, err = avif.Decode(imageReader)
		if err != nil {
			return nil, ErrFailedToDecodeAVIF
		}
	case ImageFormatICO:
		decodedImages, err := ico.Decode(imageReader)
		if err != nil {
			return nil, err
		}

		// decodedImagesの中から一番大きい画像を選択
		var largestImage *image.Image
		var largestWidth int
		var largestHeight int
		for _, di := range decodedImages {
			bounds := di.Bounds()
			if bounds.Dx() > largestWidth || bounds.Dy() > largestHeight {
				largestWidth = bounds.Dx()
				largestHeight = bounds.Dy()
				largestImage = &di
			}
		}

		if largestImage == nil {
			return nil, ErrImageFormatNotSupported
		} else {
			decodedImage = *largestImage
		}

	default:
		return nil, ErrImageFormatNotSupported
	}

	// Default quality is 75
	if targetImage.Quality == 0 {
		targetImage.Quality = 75
	}

	bounds := decodedImage.Bounds()
	resizeScale := getResizedImageScaleKeepRatio(bounds.Dx(), bounds.Dy(), targetImage.MaxWidth, targetImage.MaxHeight)

	var rgba *image.RGBA
	if resizeScale == 1.0 {
		// リサイズ不要
		rgba = image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		draw.Draw(rgba, bounds, decodedImage, bounds.Min, draw.Src)
	} else {
		newWidth := int(float64(bounds.Dx()) * resizeScale)
		newHeight := int(float64(bounds.Dy()) * resizeScale)

		rgba = image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

		// Qualityが80以上もしくはAVIFの場合はCatmullRom、それ以外はBiLinear
		if targetImage.Quality >= 80 || targetImage.Format == TargetFormatAVIF {
			draw.CatmullRom.Scale(rgba, rgba.Bounds(), decodedImage, decodedImage.Bounds(), draw.Over, nil)
		} else {
			draw.BiLinear.Scale(rgba, rgba.Bounds(), decodedImage, decodedImage.Bounds(), draw.Over, nil)
		}
	}

	if targetImage.Format == TargetFormatWebP {
		webpBytes, err := rgbaToWebP(rgba, targetImage.Quality)
		if err != nil {
			return nil, err
		} else if webpBytes == nil {
			return nil, ErrFailedToEncodeWebP
		}

		return webpBytes, nil
	} else {
		avifWriter := bytes.NewBuffer(nil)
		err := avif.Encode(avifWriter, rgba, avif.Options{
			Quality: targetImage.Quality,
		})

		if err != nil {
			return nil, ErrFailedToEncodeAVIF
		}

		avifBytes := avifWriter.Bytes()

		return &avifBytes, nil
	}
}
